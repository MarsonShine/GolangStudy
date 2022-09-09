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
*/

type ConsensusModule struct {
	mu      sync.Mutex
	id      int   // 唯一标识,CM服务器id
	peerIds []int // 集群中的其它节点id列表

	server *Server // 包含CM的服务器，用于发出RPC调用

	// 需要持久化的信息
	currentTerm int // 在所有服务器上持有的Raft状态
	votedFor    int
	log         []LogEntry

	state              CMStatus
	electionResetEvent time.Time
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
			args := RequestVoteArgs{
				Term:        savedCurrentTerm,
				CandidateId: cm.id,
			}
			var reply RequestVoteReply
			cm.dlog("sensding RequestVote to %d: %+v", peerId, args)
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

func (cm *ConsensusModule) startLeader() {

}
