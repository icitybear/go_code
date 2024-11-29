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

func getTreeRecursive(list []*Node, parentId int) []*Node {
	res := make([]*Node, 0)
	for _, v := range list {
		if v.ParentId == parentId {
			v.Children = getTreeRecursive(list, v.Id)
			res = append(res, v)
		}
	}
	return res
}

func TestDigui(t *testing.T) {
	list := []*Node{
		{4, 3, "ABA", nil},
		{3, 1, "AB", nil},
		{1, 0, "A", nil},
		{2, 1, "AA", nil},
		{5, 0, "未分配", nil},
	}
	res := getTreeRecursive(list, 0)
	bytes, _ := json.MarshalIndent(res, "", "    ")
	fmt.Printf("%s\n", bytes)

}

func getTreeIterative(list []*Node, parentId int) []*Node {
	memo := make(map[int]*Node)
	for _, v := range list {
		if _, ok := memo[v.Id]; ok {
			v.Children = memo[v.Id].Children
			memo[v.Id] = v
		} else {
			v.Children = make([]*Node, 0)
			memo[v.Id] = v
		}
		if _, ok := memo[v.ParentId]; ok {
			memo[v.ParentId].Children = append(memo[v.ParentId].Children, memo[v.Id])
		} else {
			memo[v.ParentId] = &Node{Children: []*Node{memo[v.Id]}}
		}
	}
	return memo[parentId].Children

}

func TestFor(t *testing.T) {
	list := []*Node{
		{4, 3, "ABA", nil},
		{2, 1, "AA", nil},
		{5, 3, "ABB", nil},
		{3, 1, "AB", nil},
		{1, 0, "A", nil},
		{6, 0, "未分配", nil},
	}
	res := getTreeIterative(list, 0)
	bytes, _ := json.MarshalIndent(res, "", "    ")
	fmt.Printf("%s\n", bytes)
}
