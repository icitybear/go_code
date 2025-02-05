package digui_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Node struct {
	Id       int     `json:"id"`
	ParentId int     `json:"parent_id"`
	Name     string  `json:"name"`
	Children []*Node `json:"children"`
}

// 方法一: 递归 当前节点的链路
func TestDigui(t *testing.T) {
	list := []*Node{
		{4, 3, "Grandchild2-1", nil},
		{3, 1, "Child2", nil},
		{1, 0, "Root", nil},
		{2, 1, "Child1", nil},
		{5, 0, "未分配", nil},
	}
	res := getTreeRecursive(list, 0)
	bytes, _ := json.MarshalIndent(res, "", "    ")
	fmt.Printf("%s\n", bytes)

}

func getTreeRecursive(list []*Node, nodeId int) []*Node {
	res := make([]*Node, 0)
	for _, v := range list {
		if v.ParentId == nodeId {
			v.Children = getTreeRecursive(list, v.Id)
			res = append(res, v)
		}
	}
	return res
}

// 方法2: 平铺循环 所有节点的链路
func TestFor(t *testing.T) {
	list := []*Node{
		{4, 3, "Child2-1", nil},
		{2, 1, "Child1", nil},
		{5, 3, "Child2-2", nil},
		{3, 1, "Child2", nil},
		{1, 0, "Root", nil},
		{6, 0, "未分配", nil},
	}
	res := getTreeIterative(list, 0)
	bytes, _ := json.MarshalIndent(res, "", "    ")
	fmt.Printf("%s\n", bytes)
}

func getTreeIterative(list []*Node, nodeId int) []*Node {
	memo := make(map[int]*Node)
	// 平铺循环里根据父id进行判断 归属关系
	for _, v := range list {
		if _, ok := memo[v.Id]; ok {
			v.Children = memo[v.Id].Children
			memo[v.Id] = v
		} else {
			v.Children = make([]*Node, 0)
			memo[v.Id] = v
		}
		// 包含了节点0 ParentId = 0
		if _, ok := memo[v.ParentId]; ok {
			memo[v.ParentId].Children = append(memo[v.ParentId].Children, memo[v.Id])
		} else {
			memo[v.ParentId] = &Node{Children: []*Node{memo[v.Id]}}
		}
	}

	// bytes, _ := json.MarshalIndent(memo, "", "    ")
	// fmt.Printf("%s\n", bytes)

	return memo[nodeId].Children

}

// 方法3: 集合嵌套 所有节点的链路
type Node2 struct {
	Id       int // 自增id id值随意
	Fid      int // 从0层开始算比较合理 就是特指根节点
	ParentId int // 可以没有 ParentId
	// 多了lt与rt
	LT       int
	RT       int
	Name     string
	Children []*Node2
}

func TestSet(t *testing.T) {
	list := []*Node2{
		{1, 0, 0, 1, 12, "Root", nil}, // 初始化节点是 {xxx, 0, 0, 1, 2, "Root", nil}
		{4, 1, 0, 2, 7, "Child 1", nil},
		{8, 2, 1, 3, 4, "Grandchild1-1", nil},
		{9, 3, 1, 5, 6, "Grandchild1-2", nil},
		{12, 4, 0, 8, 9, "Child 2", nil},
		{13, 5, 0, 10, 11, "Child 3", nil},
	}
	res := buildTree(list, 0)
	bytes, _ := json.MarshalIndent(res, "", "    ")
	fmt.Printf("%s\n", bytes)

}

func buildTree(list []*Node2, nodeId int) []*Node2 {

	// 创建一个映射，方便通过 ID 快速查找节点
	nodeMap := make(map[int]*Node2)
	for i := range list {
		pos := list[i].Fid
		nodeMap[pos] = list[i]
	}

	// tag: 只能根据 lt与rt的范围判断 与parentId无关
	for _, node := range list {
		for _, parentNode := range list {
			if node.LT > parentNode.LT && node.RT < parentNode.RT {
				// 判断为子节点了
				nodeMap[parentNode.Fid].Children = append(nodeMap[parentNode.Fid].Children, nodeMap[node.Fid])
				// break 不能break 不然会中断 不继续判断下个节点归属
			}
		}
	}
	bytes, _ := json.MarshalIndent(nodeMap, "", "    ")
	fmt.Printf("%s\n", bytes)

	return nodeMap[nodeId].Children
}
