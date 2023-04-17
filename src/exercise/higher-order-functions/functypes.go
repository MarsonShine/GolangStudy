package functypes

import "sort"

type State int
type States []State

// 表明参数 s 是否为目标状态
type GoalP func(s State) bool

// 成功的状态
type Successors func(s State) States

// 组合器通过将当前状态的后继者与所有其他状态组合成一个单一的状态列表来确定搜索策略。
type Combiner func(succ States, others States) States

func treeSearch(states States, goalp GoalP, succ Successors, combiner Combiner) State {
	if len(states) == 0 {
		return -1
	}
	first := states[0]
	if goalp(first) {
		return first
	} else {
		return treeSearch(combiner(succ(first), states[1:]), goalp, succ, combiner)
	}
}

// appendOthers is a Combiner function that appends others to succ.
func appendOthers(succ States, others States) States {
	return append(succ, others...)
}

// BFS
// Combiner函数，它将其他状态预加到succ中
func prependOthers(succ States, other States) States {
	return append(other, succ...)
}

func bfsTreeSearch(start State, goalp GoalP, succ Successors) State {
	return treeSearch(States{start}, goalp, succ, prependOthers)
}

func dfsTreeSearch(start State, goalp GoalP, succ Successors) State {
	return treeSearch(States{start}, goalp, succ, appendOthers)
}

func binaryTree(s State) States {
	return []State{s * 2, s*2 + 1}
}

func finiteBinaryTree(n State) Successors {
	return func(s State) States {
		return filter(binaryTree(s), func(item State) bool { return item <= n })
	}
}

func filter[T any](s []T, pred func(item T) bool) []T {
	var result []T
	for _, item := range s {
		if pred(item) {
			result = append(result, item)
		}
	}
	return result
}

// stateIs返回一个检查状态是否与n相等的GoalP。
func stateIs(n State) GoalP {
	return func(s State) bool { return n == s }
}

type CostFunc func(s State) int

// 根据 cost 排序
func sorter(cost CostFunc) Combiner {
	return func(succ States, others States) States {
		all := append(succ, others...)
		sort.Slice(all, func(i, j int) bool {
			return cost(all[i]) < cost(all[j])
		})
		return all
	}
}

func bestCostTreeSearch(start State, goalp GoalP, succ Successors, cost CostFunc) State {
	return treeSearch(States{start}, goalp, succ, sorter(cost))
}

// costDiffTarget 使用数字从`n`的距离作为一个cost参数来创建 CostFunc
func costDiffTarget(n State) CostFunc {
	return func(s State) int {
		delta := int(s) - int(n)
		if delta < 0 {
			return -delta
		} else {
			return delta
		}
	}
}
