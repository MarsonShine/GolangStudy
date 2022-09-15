package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

/*
raft 论文地址：https://raft.github.io/raft.pdf

共识算法（Consensus Module）是raft算法的核心；它完全抽象了与集群中其他副本的网络和连接的细节

概念介绍：
peer：在raft中，集群中的每个副本都可以称其它的副本为”对等体“，而peerIds就是除自己外的其它副本的id列表
server: 它使ConsensusModule能够向对等体发送消息

CM状态机有三种状态：Leader，Follower，Candidate；三种状态的转换图详见：https://eli.thegreenplace.net/images/2020/raft-highlevel-state-machine.png
Terms（任期）：从节点每次参与选举（election）都会有自己的一个任期数，Raft算法确保给定的任期有一个leader

Election Timer（选举计时器）：Raft算法其中一个关键组成部分就是选举计时器（也叫心跳）；follower节点会一直运行计时器，每得到当前leader的响应就会重新启动它。当leader无法在规定的时间做出响应，从节点就认为leader断开连接，于是就开始进行选举（从而切换到了Candidate状态）
Q：这个时候是不是所有的follower节点会同时成为candidate？
A：并不是，由于这个选举计时器的启动是随机的，所以这大大减少了follower节点同时称为candidate的可能性。但是即使存在这种同时称为候选人，也只有一个节点会在这个任期内称为leader。在极少数情况下选举出现分裂，以至于在这次的任期内没能选出leader。但也会在下一次term中选出leader。当然也可能存在无限循环选举的可能，这种情况发生的概率极其低
Q：如果集群从节点因为网络分区，那么它会因为没有leader的消息而开始选举吗？
A：是的，它会进行选择。但是它不会得到任期的，因为它不能连接其它的peer节点（或某一部分），因此得不到选票。它可能会一直在这个candidate状态下重复自我选举，直到它重新连接到集群上。

深入内部peer节点之间的rpc请求
rafe有两种类别的rpc请求：
1. 请求投票（RequestVotes）：这个只在candidate状态使用；candidate节点像其它peer节点发出投票请求，其返回结果包含了是否得到了选票。
2. 日志请求，追加条目（AppendEntries）：只有在leader下使用；leader会将复制日志广播给其它从节点，即使没有新的日志，也会定期广播其它节点保持心跳
所以我们可以得知，从节点是不会主动发出rpc的，只有在它们自己的选举计时器在规定的时间内没有得到leader的响应，就会从follower转变成candidate并发送RV。

如果距离上一次的“选举重置事件”已经有一段时间了，这个peer节点就会开始选举，并成为candidate节点。
这个事件是什么？它可以是终止选举的任何事情--例如收到了有效的心跳，或者投票给了另一个candidate节点。

成为了候选者（Candidate）之后呢？
我们在代码中已经看到，一旦足够长的时间过去了，follower没有听到leader心跳或其他候选者的消息，选举就开始了。在查看代码之前，让我们讨论一下进行选举需要做的事情:
1. 将状态转换为candidate，并将当前任期数自增（表示新一轮投票）
2. 发送RV给所有的peers节点，并请求为自己投票
3. 等待这些RPC的响应，并统计投票数，如果数量足够了就变成leader

成为leader后，需要向其它节点同步信息，发送RPC请求

状态转移和协程
总结一下cm状态的转移过程和对应转移所处的goroutine：
Follower：CM被初始化为一个follower，每次调用becomeFollower时都会新起goroutine运行runElectionTimer。注意，在短时间内允许有多个同时运行。假设一个follower从更高任期的leader接收到一个RV，就会触发becomeFollower调用，就会启动新的goroutine定时器。但是旧的routine计时器会因为任期不是最新从而退出。

Candidate：候选人也有并行的选举goroutine，但除此之外，它还有许多goroutine来发送rpc。它同样也有follower那样的机制，根据任期大小来停止旧的选举。请记住，RPC goroutines可能需要很长的时间才能完成，所以如果它们发现在RPC调用返回时已经过期，就必须安静地退出。

Leader：leader 不会发出选举goroutine，但是会以每50毫秒来发送心跳


客户端交互部分
客户端是如何发现leader？
客户端会通过raft来复制一连串的command，这些命令作为通用状态机的输入，这些命令一般会经过如下阶段：
1. 命令由客户端提交给leader。在raft对等体集群中，一个命令通常只提交给一个对等体（peer）
2. leader复制这个命令给其它followers
3. 一旦leader对命令被充分复制了即满足提交条件了（即大多数集群对等体在日志中都记录这个命令），这个命令就被提交并并通知到其它客户端有新提交。


*/

type ConsensusModule struct {
	mu      sync.Mutex
	id      int   // 唯一标识,CM服务器id
	peerIds []int // 集群中的其它节点id列表

	server *Server // 包含CM的服务器，用于发出RPC调用

	// 用于报告已提交的日志条目的通道，由客户端构建的过程中传入的
	commitChan chan<- CommitEntry

	// newCommitReadyChan是一个内部通知通道，由向日志提交新条目的goroutine使用，以通知这些条目可以在commitChan上发送。
	newCommitReadyChan chan struct{}

	// 需要持久化的信息
	currentTerm int // 在所有服务器上持有的Raft状态
	votedFor    int
	log         []LogEntry

	// 所有节点服务器的易失性状态
	commitIndex        int
	lastApplied        int
	state              CMStatus
	electionResetEvent time.Time

	// leader的易失性状态
	nextIndex  map[int]int
	matchIndex map[int]int
}

// raft 论文中的图2介绍了日志结构
type AppendEntriesArgs struct {
	Term     int
	LeaderId int

	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}
type AppendEntriesReply struct {
	Term    int
	Success bool
}

type LogEntry struct {
	Command interface{}
	Term    int
}

type CMStatus int

const (
	Follower CMStatus = iota
	Candidate
	Leader
	Dead
)

type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// 要提交得条目，每个提交条目都会通知客户端，在提交得命名上达成共识，最终应用到客户端上。
type CommitEntry struct {
	Command interface{}
	Index   int
	Term    int
}

/*
首先选定一个超时时间，如论文建议150-300毫秒，同样在访问cm成员对象时需要上锁。这点非常重要，对状态的更新是同步操作
循环体是每10毫秒执行一次。
检查cm的状态不是candidate以及不是follower看起来有点奇怪，不通过runElectionTimer启动的选举，这个节点能突然成为leader么？
如果不是预期的状态就会退出计时器。
*/
func (cm *ConsensusModule) runElectionTimer() {
	timeoutDuration := cm.electionTimeout()
	cm.mu.Lock()
	termStarted := cm.currentTerm
	cm.mu.Unlock()
	cm.dlog("election timer started (%v), term=%d", timeoutDuration, termStarted)

	// 循环结束条件：
	// 1.不在需要选举计时器 2. 选举计时器超时，以及CM成为了candidate
	// 在从节点中，都会在后台持续运行，一直到CM死亡
	timer := time.NewTicker(10 * time.Millisecond)
	defer timer.Stop()
	for {
		<-timer.C
		cm.mu.Lock()
		if cm.state != Candidate && cm.state != Follower {
			cm.dlog("in election timer state=%s, bailing out", cm.state)
			cm.mu.Unlock()
			return
		}
		if termStarted != cm.currentTerm {
			cm.dlog("in election timer term changed from %d to %d, bailing out", termStarted, cm.currentTerm)
			cm.mu.Unlock()
			return
		}

		// 如果在超时时间内没有接收到leader的心跳，或者没有投票，那么就开始选举
		if elapsed := time.Since(cm.electionResetEvent); elapsed >= timeoutDuration {
			cm.startElection()
			cm.mu.Unlock()
			return
		}
		cm.mu.Unlock()
	}
}

func (cm *ConsensusModule) electionTimeout() time.Duration {
	// 通过环境变量读取配置的超时时间
	if len(os.Getenv("RAFT_FORCE_MORE_REELECTION")) > 0 && rand.Intn(3) == 0 {
		return time.Duration(150) * time.Millisecond
	} else {
		return time.Duration(150+rand.Intn(150)) * time.Millisecond
	}
}

func (cm *ConsensusModule) dlog(format string, args ...interface{}) {
	format = fmt.Sprintf("[%d] ", cm.id) + format
	log.Printf(format, args...)
}

/*
候选者通过给自己投票开始的，即设置votesReceived=1，并将投票人设为自己：votedFor = cm.id
然后向其他peers节点发送RV，并等待所有响应结果；
比较响应的结果，如果获得超半数的投票，则成为leader。如果是其他节点成为leader，则将自身candidate状态变更为follower；
最后再开启新的选择计时器
*/
func (cm *ConsensusModule) startElection() {
	cm.state = Candidate
	cm.currentTerm++
	savedCurrentTerm := cm.currentTerm
	cm.electionResetEvent = time.Now()
	cm.votedFor = cm.id
	cm.dlog("becomes Candidate (currentTerm=%d); log=%v", savedCurrentTerm, cm.log)

	votesReceived := 1

	// 发送RV给其他节点
	for _, peerId := range cm.peerIds {
		go func(peerId int) {
			cm.mu.Lock()
			savedLastLogIndex, savedLastLogTerm := cm.lastLogIndexAndTerm()
			cm.mu.Unlock()
			args := RequestVoteArgs{
				Term:         savedCurrentTerm,
				CandidateId:  cm.id,
				LastLogIndex: savedLastLogIndex,
				LastLogTerm:  savedLastLogTerm,
			}
			cm.dlog("sending RequestVote to %d: %+v", peerId, args)

			var reply RequestVoteReply
			if err := cm.server.Call(peerId, "ConsensusModule.RequestVote", args, &reply); err != nil {
				cm.mu.Lock()
				defer cm.mu.Unlock()
				cm.dlog("received RequestVoteReply %+v", reply)

				if cm.state != Candidate {
					cm.dlog("while waiting for reply, state = %v", cm.state)
					return
				}

				if reply.Term > savedCurrentTerm { // term比该发出RV的cm的term大，说明此次选举无效，过期。并已经有leader节点了
					cm.dlog("term out of date in RequestVoteReply")
					cm.becomeFollower(reply.Term)
					return
				} else if reply.Term == savedCurrentTerm {
					if reply.VoteGranted {
						votesReceived++
						if votesReceived*2 > len(cm.peerIds)+1 { // 获得大多数选举投票
							cm.dlog("wins election with %d votes", votesReceived)
							cm.startLeader()
							return
						}
					}
				}
			}
		}(peerId)
	}

	// 运行另一个选举计时器，以防这次选举不成功
	go cm.runElectionTimer()
}

func (cm *ConsensusModule) lastLogIndexAndTerm() (int, int) {
	if len(cm.log) > 0 {
		lastIndex := len(cm.log) - 1
		return lastIndex, cm.log[lastIndex].Term
	} else {
		return -1, -1 // 哨兵机制
	}
}

// 将该cm状态变更为follower并重置state，该操作必须在锁内
func (cm *ConsensusModule) becomeFollower(term int) {
	cm.dlog("becomes Follower with term=%d; log=%v", term, cm.log)
	cm.state = Follower
	cm.currentTerm = term
	cm.votedFor = -1
	cm.electionResetEvent = time.Now()
	// 继续开启下一轮选举计时器
	go cm.runElectionTimer()
}

/*
赢得大多数节点的选票即可称为leader
每50ms发送一个心跳请求
*/
func (cm *ConsensusModule) startLeader() {
	cm.state = Leader
	cm.dlog("becomes Leader; term=%d, log=%v", cm.currentTerm, cm.log)

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		// 作为leader，需要一直发送心跳
		for {
			cm.leaderSendHeartbeats()
			<-ticker.C

			cm.mu.Lock()
			if cm.state != Leader {
				cm.mu.Unlock()
				return
			}
			cm.mu.Unlock()
		}
	}()
}

func (cm *ConsensusModule) leaderSendHeartbeats() {
	cm.mu.Lock()
	savedCurrentTerm := cm.currentTerm
	cm.mu.Unlock()

	for _, peerId := range cm.peerIds {
		go func(peerId int) {
			cm.mu.Lock()
			ni := cm.nextIndex[peerId]
			prevLogIndex := ni - 1
			prevLogTerm := -1
			if prevLogIndex >= 0 {
				prevLogTerm = cm.log[prevLogIndex].Term
			}
			entries := cm.log[ni:]

			args := AppendEntriesArgs{
				Term:         savedCurrentTerm,
				LeaderId:     cm.id,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: cm.commitIndex,
			}
			cm.mu.Unlock()
			cm.dlog("sending AppendEntries to %v: ni=%d, args=%+v", peerId, ni, args)
			// AE响应有一个success字段，告诉leader follower是否看到了prevLogIndex和prevLogTerm的匹配。基于这个字段，leader为这个follower更新nextIndex。
			var reply AppendEntriesReply
			if err := cm.server.Call(peerId, "ConsensusModule.AppendEntries", args, &reply); err == nil {
				cm.mu.Lock()
				defer cm.mu.Unlock()
				if reply.Term > savedCurrentTerm {
					cm.dlog("term out of date in heartbeat reply")
					cm.becomeFollower(reply.Term)
					return
				}

				if cm.state == Leader && savedCurrentTerm == reply.Term {
					if reply.Success {
						cm.nextIndex[peerId] = ni + len(entries)
						cm.matchIndex[peerId] = cm.nextIndex[peerId] - 1
						cm.dlog("AppendEntries reply from %d success: nextIndex := %v, matchIndex := %v", peerId, cm.nextIndex, cm.matchIndex)

						savedCommitIndex := cm.commitIndex
						for i := cm.commitIndex; i < len(cm.log); i++ {
							if cm.log[i].Term == cm.currentTerm {
								matchCount := 1
								for _, peerId := range cm.peerIds {
									if cm.matchIndex[peerId] >= i {
										matchCount++
									}
								}
								// commitIndex只有在大多数follower复制成功日志索引，才会更新为这个索引位置
								if matchCount*2 > len(cm.peerIds)+1 {
									cm.commitIndex = i
								}
							}
						}
						if cm.commitIndex != savedCommitIndex {
							cm.dlog("leader sets commitIndex := %d", cm.commitIndex)
							cm.newCommitReadyChan <- struct{}{}
						}
					} else {
						cm.nextIndex[peerId] = ni - 1
						cm.dlog("AppendEntries reply from %d !success: nextIndex := %d", peerId, ni-1)
					}
				}
			}
		}(peerId)
	}
}

func (cm *ConsensusModule) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.state == Dead {
		return nil
	}
	lastLogIndex, lastLogTerm := cm.lastLogIndexAndTerm()
	cm.dlog("RequestVote: %+v [currentTerm=%d, votedFor=%d, log index/term=(%d, %d)]", args, cm.currentTerm, cm.votedFor, lastLogIndex, lastLogTerm)

	if args.Term > cm.currentTerm {
		// 请求参数的任期比该cm大，说明该投票已过期
		cm.dlog("... term out of date in RequestVote")
		cm.becomeFollower(args.Term)
	}
	// 调用者的任期和请求投票的任期一致并且该调用者还没有投票给其它节点请求，则投票成功
	if cm.currentTerm == args.Term && (cm.votedFor == -1 || cm.votedFor == args.CandidateId) &&
		(args.LastLogTerm > lastLogTerm && args.LastLogIndex >= lastLogIndex) {
		reply.VoteGranted = true
		cm.votedFor = args.CandidateId
		cm.electionResetEvent = time.Now()
	} else {
		reply.VoteGranted = false
	}
	reply.Term = cm.currentTerm
	cm.dlog("... RequestVote reply: %+v", reply)
	return nil

}

/*
这段代码严格遵循论文中图2的算法（AppendEntries的Receive实现部分），而且注释得很好。
注意，当代码注意到leader的LeaderCommit大于自己的自己的cm.commitIndex时，就会在cm.newCommitReadyChan上发送；
这时follower意识到leader认为有额外的条目要提交。
当leader通过AppendEntries发送日志条目时，会发生以下情况：
1. 一个follower在它的日志中追加新的条目并将 success=true 返回给leader
2. 作为结果，leader 会根据 follwer 更新 matchIndex。当足够数量的 follower 在 nextIndex 都有matchIndex时
leader就会更新 commitIndex 并在下一个AE中发送给所有的follwers（在 leaderCommit 字段中）
3. 当 followers 接收到一个超出他们之前所见的新的 leaderCommit 时，它们知道有新的日志条目已经提交，它们可以在通道上将其发送给他们的客户
Q: 提交一个新命令需要多少次RPC往返?
A: 两次；一次是 leader 发送下一个日志条目给 followers 以及 followers 返回 ack；二次是将发送更新 commitIndex 给 followers，然后 followers 将这些条目标记为已提交，并将他们发送到提交通道。
*/
func (cm *ConsensusModule) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.state == Dead {
		return nil
	}
	cm.dlog("AppendEntries: %+v", args)

	if args.Term > cm.currentTerm {
		cm.dlog("... term out of date in AppendEntries")
		cm.becomeFollower(args.Term)
	}

	reply.Success = false
	/*
		这里在什么情况下会存在一个leader(该cm)节点会成为另一个leader的follower节点？
		Raft保证在任何给定的任期内只存在一个领导者。如果您仔细地遵循RequestVote的逻辑和startElection中发送rv的代码，您将看到集群中不可能在同一任期内存在两个leader。这个条件对于发现另一个同伴在这届选举中获胜的候选人来说是很重要的。
		这里应该指的是短时间leader失去连接，导致脑裂，从而出现的问题
	*/
	if args.Term == cm.currentTerm {
		if cm.state != Follower {
			cm.becomeFollower(args.Term)
		}
		cm.electionResetEvent = time.Now()
		// 日志是否包含一个PrevLogIndex位置的条目，其任期是否匹配？请注意，在PrevLogIndex=-1的极端情况下，这肯定是真的。
		if args.PrevLogIndex == -1 ||
			(args.PrevLogIndex < len(cm.log)) && args.PrevLogTerm == cm.log[args.PrevLogIndex].Term {
			reply.Success = true
			// 找到插入点，从PrevLogIndex+1开始的现有日志和RPC发送的新条目之间存在任期不匹配的地方。
			logInsertIndex := args.PrevLogIndex + 1
			newEntriesIndex := 0

			for {
				if logInsertIndex >= len(cm.log) || newEntriesIndex >= len(args.Entries) {
					break
				}
				if cm.log[logInsertIndex].Term != args.Entries[newEntriesIndex].Term {
					break
				}
				logInsertIndex++
				newEntriesIndex++
			}
			// 循环最后部分：
			// - logInsertIndex 指向日志索引末尾位置，或者是在该索引上的任期与leader上的不匹配
			// - newEntriesIndex 指向Entries的末尾位置，或者指向与相应的日志条目的任期不匹配的索引
			if newEntriesIndex < len(args.Entries) {
				cm.dlog("... inserting entries %v from index %d", args.Entries[newEntriesIndex:], logInsertIndex)
				cm.log = append(cm.log[:logInsertIndex], args.Entries[newEntriesIndex:]...)
				cm.dlog("... log is now: %v", cm.log)
			}

			// 赋值commitIndex
			if args.LeaderCommit > cm.commitIndex {
				cm.commitIndex = intMin(args.LeaderCommit, len(cm.log)-1)
				cm.dlog("... setting commitIndex=%d", cm.commitIndex)
				cm.newCommitReadyChan <- struct{}{}
			}
		}
		reply.Success = true
	}
	reply.Term = cm.currentTerm
	cm.dlog("AppendEntries reply: %+v", *reply)
	return nil
}

func (cm *ConsensusModule) Submit(command interface{}) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.dlog("Submit received by %v: %v", cm.state, command)
	if cm.state == Leader {
		cm.log = append(cm.log, LogEntry{Command: command, Term: cm.currentTerm})
		cm.dlog("... log=%v", cm.log)
		return true
	}
	return false
}

// 该方法在goroutine启动时就会运行
// 该方法更新lastApplied状态变量，以了解哪些条目已经发送到客户端，并只发送新的条目。
func (cm *ConsensusModule) commitChanSender() {
	for range cm.newCommitReadyChan {
		cm.mu.Lock()
		savedTerm := cm.currentTerm
		savedLastApplied := cm.lastApplied
		var entries []LogEntry
		if cm.commitIndex > cm.lastApplied {
			entries = cm.log[cm.lastApplied+1 : cm.commitIndex]
			cm.lastApplied = cm.commitIndex
		}
		cm.mu.Unlock()
		cm.dlog("commitChanSender entries=%v, savedLastApplied=%d", entries, savedLastApplied)

		for i, entry := range entries {
			cm.commitChan <- CommitEntry{
				Command: entry.Command,
				Index:   savedLastApplied + i + 1,
				Term:    savedTerm,
			}
		}
	}
	cm.dlog("commiteChanSender done")
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
