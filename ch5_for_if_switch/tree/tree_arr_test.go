package tree_test

import (
	"fmt"
	"testing"
)

/* 数组表示下的二叉树类 */
type arrayBinaryTree struct {
	tree []any
}

/* 构造方法 */
func newArrayBinaryTree(arr []any) *arrayBinaryTree {
	return &arrayBinaryTree{
		tree: arr,
	}
}

/* 列表容量 */
func (abt *arrayBinaryTree) size() int {
	return len(abt.tree)
}

/* 获取索引为 i 节点的值 */
func (abt *arrayBinaryTree) val(i int) any {
	// 若索引越界，则返回 null ，代表空位
	if i < 0 || i >= abt.size() {
		return nil
	}
	return abt.tree[i]
}

/* 获取索引为 i 节点的左子节点的索引 */
func (abt *arrayBinaryTree) left(i int) int {
	return 2*i + 1
}

/* 获取索引为 i 节点的右子节点的索引 */
func (abt *arrayBinaryTree) right(i int) int {
	return 2*i + 2
}

/* 获取索引为 i 节点的父节点的索引 */
func (abt *arrayBinaryTree) parent(i int) int {
	return (i - 1) / 2
}

/* 层序遍历 */
func (abt *arrayBinaryTree) levelOrder() []any {
	var res []any
	// 直接遍历数组
	for i := 0; i < abt.size(); i++ {
		if abt.val(i) != nil {
			res = append(res, abt.val(i))
		}
	}
	return res
}

/* 深度优先遍历 */
func (abt *arrayBinaryTree) dfs(i int, order string, res *[]any) {
	// 若为空位，则返回
	if abt.val(i) == nil {
		return
	}
	// 前序遍历
	if order == "pre" {
		*res = append(*res, abt.val(i))
	}
	abt.dfs(abt.left(i), order, res)
	// 中序遍历
	if order == "in" {
		*res = append(*res, abt.val(i))
	}
	abt.dfs(abt.right(i), order, res)
	// 后序遍历
	if order == "post" {
		*res = append(*res, abt.val(i))
	}
}

/* 前序遍历 */
func (abt *arrayBinaryTree) preOrder() []any {
	var res []any
	abt.dfs(0, "pre", &res)
	return res
}

/* 中序遍历 */
func (abt *arrayBinaryTree) inOrder() []any {
	var res []any
	abt.dfs(0, "in", &res)
	return res
}

/* 后序遍历 */
func (abt *arrayBinaryTree) postOrder() []any {
	var res []any
	abt.dfs(0, "post", &res)
	return res
}

func TestArr(t *testing.T) {
	/* 二叉树的数组表示 */
	// 使用 any 类型的切片, 就可以使用 nil 来标记空位
	tree := []any{1, 2, 3, 4, nil, 6, 7, 8, 9, nil, nil, 12, nil, nil, 15}
	bst := newArrayBinaryTree(tree)

	fmt.Println(bst.val(3)) // 4
	fmt.Println(bst.val(4)) // <nil>
	fmt.Println("初始化：层序遍历:")
	arr := bst.levelOrder()
	fmt.Println(arr) // [1 2 3 4 6 7 8 9 12 15]

	fmt.Println("初始化：前序遍历:")
	arr = bst.preOrder()
	fmt.Println(arr) // [1 2 4 8 9 3 6 12 7 15]
}
