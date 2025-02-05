package tree_test

import (
	"container/list"
	"fmt"
	"testing"
)

// 定义节点结构体 链表
type Node struct {
	value int // 测试存放整形数据
	left  *Node
	right *Node
}

// 定义二叉排序树结构体
type BinarySortTree struct {
	root *Node
}

// 创建一个新的节点
func NewNode(value int) *Node {
	return &Node{value: value, left: nil, right: nil}
}

// 创建一个空的二叉排序树
func NewBinarySortTree() *BinarySortTree {
	return &BinarySortTree{root: nil}
}

// 清空树
func (bst *BinarySortTree) MakeEmpty() {
	bst.root = nil
}

// 判断树是否为空
func (bst *BinarySortTree) IsEmpty() bool {
	return bst.root == nil
}

// 查找最小值 递归
func (bst *BinarySortTree) FindMin(node *Node) *Node {
	if node == nil {
		return nil
	} else if node.left == nil {
		return node
	} else {
		return bst.FindMin(node.left)
	}
}

// 查找最大值 递归
func (bst *BinarySortTree) FindMax(node *Node) *Node {
	if node == nil {
		return nil
	} else if node.right == nil {
		return node
	} else {
		return bst.FindMax(node.right)
	}
}

// 判断值是否存在 遍历链表
func (bst *BinarySortTree) IsContains(x int) bool {
	current := bst.root
	for current != nil {
		if x < current.value {
			current = current.left
		} else if x > current.value {
			current = current.right
		} else {
			return true
		}
	}
	return false
}

// 插入值
func (bst *BinarySortTree) Insert(x int) *Node {
	if bst.root == nil {
		bst.root = NewNode(x)
		return bst.root
	}

	current := bst.root
	for current != nil {
		if x < current.value {
			if current.left == nil {
				current.left = NewNode(x)
				return current.left
			} else {
				current = current.left
			}
		} else if x > current.value {
			if current.right == nil {
				current.right = NewNode(x)
				return current.right
			} else {
				current = current.right
			}
		} else {
			// 如果值已经存在，直接返回当前节点
			return current
		}
	}
	return nil
}

// 删除值 递归处理含有双节点的
func (bst *BinarySortTree) Remove(x int, node *Node) *Node {
	if node == nil {
		return nil
	}

	if x < node.value {
		node.left = bst.Remove(x, node.left)
	} else if x > node.value {
		node.right = bst.Remove(x, node.right)
	} else if node.left != nil && node.right != nil {
		// 找到右子树的最小值替换当前节点
		node.value = bst.FindMin(node.right).value
		node.right = bst.Remove(node.value, node.right)
	} else {
		if node.left == nil && node.right == nil {
			node = nil
		} else if node.right != nil {
			node = node.right
		} else if node.left != nil {
			node = node.left
		}
		return node
	}
	return node
}

func TestBinarySortTree(t *testing.T) {
	// 测试二叉排序树
	bst := NewBinarySortTree()
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(2)
	bst.Insert(4)
	bst.Insert(6)
	bst.Insert(8)

	fmt.Println("初始化：中序遍历:")
	bst.InOrderTraversal(bst.root) // 2 3 4 5 6 7 8
	fmt.Println()

	fmt.Println("4是否在树里?", bst.IsContains(4))
	fmt.Println("9是否在树里?", bst.IsContains(9))

	fmt.Println("Min value:", bst.FindMin(bst.root).value)
	fmt.Println("Max value:", bst.FindMax(bst.root).value)

	bst.Remove(3, bst.root)
	fmt.Println("删除节点3后，中序遍历:")
	bst.InOrderTraversal(bst.root)
	fmt.Println()
}

// 层序遍历-使用数组
func (bst *BinarySortTree) LevelOrderTraversal() {
	if bst.root == nil {
		return
	}

	// 使用数组或队列实现层序遍历
	queue := []*Node{bst.root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:] // 出队
		fmt.Printf("%d ", current.value)

		// 左子节点入队
		if current.left != nil {
			queue = append(queue, current.left)
		}
		// 右子节点入队
		if current.right != nil {
			queue = append(queue, current.right)
		}
	}
}

// 层序遍历-使用队列 int可以使用任何类型any替换
func levelOrder(root *Node) []int {
	// 初始化队列，加入根节点
	queue := list.New()
	queue.PushBack(root)
	// 初始化一个切片，用于保存遍历序列
	nums := make([]int, 0)
	for queue.Len() > 0 {
		// 队列出队
		node := queue.Remove(queue.Front()).(*Node)
		// 保存节点值
		nums = append(nums, node.value)
		if node.left != nil {
			// 左子节点入队
			queue.PushBack(node.left)
		}
		if node.right != nil {
			// 右子节点入队
			queue.PushBack(node.right)
		}
	}
	return nums
}

func TestLevel(t *testing.T) {
	bst := NewBinarySortTree()
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(2)
	bst.Insert(4)
	bst.Insert(6)
	bst.Insert(8)

	fmt.Println("初始化：层序遍历:")
	bst.LevelOrderTraversal() // 5 3 7 2 4 6 8
	fmt.Println()

	fmt.Println("初始化：层序遍历2队列:")
	arr := levelOrder(bst.root)
	fmt.Println(arr) // [5 3 7 2 4 6 8]
}

func TestOrderDigui(t *testing.T) {
	bst := NewBinarySortTree()
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(2)
	bst.Insert(4)
	bst.Insert(6)
	bst.Insert(8)
	fmt.Println("初始化：前序遍历:")
	bst.PreOrderTraversal(bst.root) // 5 3 2 4 7 6 8
	fmt.Println()

	fmt.Println("初始化：中序遍历:")
	bst.InOrderTraversal(bst.root) // 2 3 4 5 6 7 8
	fmt.Println()

	fmt.Println("初始化：后序遍历:")
	bst.PostfixOrderTraversal(bst.root) // 2 4 3 6 8 7 5
	fmt.Println()
}

/* 前序遍历 */
func (bst *BinarySortTree) PreOrderTraversal(node *Node) {
	if node == nil {
		return
	}
	// 访问优先级：根节点 -> 左子树 -> 右子树
	fmt.Printf("%d ", node.value)
	bst.PreOrderTraversal(node.left)
	bst.PreOrderTraversal(node.right)
}

/* 中序遍历 */
func (bst *BinarySortTree) InOrderTraversal(node *Node) {
	if node != nil {
		// 访问优先级：左子树 -> 根节点 -> 右子树
		bst.InOrderTraversal(node.left)
		fmt.Printf("%d ", node.value)
		bst.InOrderTraversal(node.right)
	}
}

// func inOrder(node *Node) {
// 	if node == nil {
// 		return
// 	}
// 	// 访问优先级：左子树 -> 根节点 -> 右子树
// 	inOrder(node.left)
// 	nums = append(nums, node.value)
// 	inOrder(node.right)
// }

/* 后序遍历 */
func (bst *BinarySortTree) PostfixOrderTraversal(node *Node) {
	if node == nil {
		return
	}
	// 访问优先级：左子树 -> 右子树 -> 根节点
	bst.PostfixOrderTraversal(node.left)
	bst.PostfixOrderTraversal(node.right)
	fmt.Printf("%d ", node.value)
}

// 前序遍历（非递归）
func (bst *BinarySortTree) PreOrderStack() {
	if bst.root == nil {
		return
	}

	stack := []*Node{bst.root}
	for len(stack) > 0 {
		// 弹出栈顶节点
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		fmt.Printf("%d ", current.value)

		// 右子节点先入栈，左子节点后入栈
		if current.right != nil {
			stack = append(stack, current.right)
		}
		if current.left != nil {
			stack = append(stack, current.left)
		}
	}
	fmt.Println()
}

// 中序遍历（非递归）
func (bst *BinarySortTree) InOrderStack() {
	if bst.root == nil {
		return
	}

	stack := []*Node{}
	current := bst.root
	for current != nil || len(stack) > 0 {
		// 将左子节点全部入栈
		for current != nil {
			stack = append(stack, current)
			current = current.left
		}

		// 弹出栈顶节点并访问
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		fmt.Printf("%d ", current.value)

		// 转向右子节点
		current = current.right
	}
	fmt.Println()
}

// 后序遍历（非递归）
func (bst *BinarySortTree) PostfixOrderStack() {
	if bst.root == nil {
		return
	}

	stack := []*Node{bst.root}
	var prev *Node // 记录上一个访问的节点

	for len(stack) > 0 {
		current := stack[len(stack)-1]

		// 如果当前节点是叶子节点，或者已经访问过其子节点
		if (current.left == nil && current.right == nil) || (prev != nil && (prev == current.left || prev == current.right)) {
			fmt.Printf("%d ", current.value)
			stack = stack[:len(stack)-1]
			prev = current
		} else {
			// 右子节点先入栈，左子节点后入栈
			if current.right != nil {
				stack = append(stack, current.right)
			}
			if current.left != nil {
				stack = append(stack, current.left)
			}
		}
	}
	fmt.Println()
}

func TestOrderStack(t *testing.T) {
	bst := NewBinarySortTree()
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(2)
	bst.Insert(4)
	bst.Insert(6)
	bst.Insert(8)

	fmt.Println("初始化：前序遍历-非递归:")
	bst.PreOrderTraversal(bst.root) // 5 3 2 4 7 6 8
	fmt.Println()

	fmt.Println("初始化：中序遍历-非递归:")
	bst.InOrderTraversal(bst.root) // 2 3 4 5 6 7 8
	fmt.Println()

	fmt.Println("初始化：后序遍历-非递归:")
	bst.PostfixOrderTraversal(bst.root) // 2 4 3 6 8 7 5
	fmt.Println()
}
