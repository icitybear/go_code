package pre_tree_test

import (
	"fmt"
	"strings"
	"testing"
)

type radixNode struct {
	path     string       // 当前节点的相对路径
	fullPath string       // 完整路径
	indices  string       // 每个 indice 字符对应一个孩子节点的 path 首字母
	children []*radixNode // 后继节点切片
	end      bool         // 是否有路径以当前节点为终点
	passCnt  int          // 记录有多少路径途径当前节点
}

//radix tree 和 trie 所不同的是，root 节点是能够存储实际的相对路径 path 的.

type Radix struct {
	root *radixNode
}

func NewRadix() *Radix {
	return &Radix{
		root: &radixNode{},
	}
}

func (r *Radix) Insert(word string) {
	// 不重复插入
	if r.Search(word) {
		return
	}
	r.root.insert(word)
}

// 插入节点流程
func (rn *radixNode) insert(word string) {
	fullWord := word

	// 如果当前节点为 root，此之前没有注册过子节点，则直接插入并返回
	if rn.path == "" && len(rn.children) == 0 {
		rn.insertWord(word, word)
		return
	}

walk:
	for {
		// 获取到 word 和当前节点 path 的公共前缀长度 i 进行分类处理
		i := commonPrefixLen(word, rn.path)
		// 1. 只要公共前缀大于 0
		if i > 0 {
			rn.passCnt++ // 说明 word 必然经过当前节点，需要对 passCnt 计数器加 1
		}

		// 2. 公共前缀小于当前节点的相对路径，要将当前节点拆分为公共前缀部分 + 后继剩余部分两个节点
		if i < len(rn.path) {
			// 需要进行节点切割
			child := radixNode{
				// 进行相对路径切分
				path: rn.path[i:],
				// 继承完整路径
				fullPath: rn.fullPath,
				// 当前节点的后继节点进行委托
				children: rn.children,
				indices:  rn.indices,
				end:      rn.end,
				// 传承给孩子节点时，需要把之前累加上的 passCnt 计数扣除
				passCnt: rn.passCnt - 1,
			}

			// 续接上孩子节点
			rn.indices = string(rn.path[i])
			rn.children = []*radixNode{&child}
			// 调整原节点的 full path
			rn.fullPath = rn.fullPath[:len(rn.fullPath)-(len(rn.path)-i)]
			// 调整原节点的 path
			rn.path = rn.path[:i]
			// 原节点是新拆分出来的，目前不可能有单词以该节点结尾
			rn.end = false
		}

		// 3. 公共前缀小于插入 word 的长度，则需要继续检查，word 和后继节点是否还有公共前缀，如果有的话，则递归对后继节点执行相同流程
		if i < len(word) {
			// 对 word 扣除公共前缀部分
			word = word[i:]
			// 获取 word 剩余部分的首字母
			c := word[0]
			for i := 0; i < len(rn.indices); i++ {
				// 如果与后继节点还有公共前缀，则将 rn 指向子节点，然后递归执行流程
				if rn.indices[i] == c {
					rn = rn.children[i]
					continue walk
				}
			}

			// 到了这里，意味着 word 剩余部分与后继节点没有公共前缀了
			// 则直接将 word 包装成一个新的节点，插入到当前节点的子节点列表 children 当中
			rn.indices += string(c)
			child := radixNode{}
			child.insertWord(word, fullWord) // 传入相对路径和完整路径，补充一个新生成的节点信息
			rn.children = append(rn.children, &child)
			return
		}

		// 倘若公共前缀恰好是 path，需要将 end 置为 true
		rn.end = true
		return
	}
}

// 求取两个单词的公共前缀
func commonPrefixLen(wordA, wordB string) int {
	var move int
	for move < len(wordA) && move < len(wordB) && wordA[move] == wordB[move] {
		move++
	}
	return move
}

// 传入相对路径和完整路径，补充一个新生成的节点信息
func (rn *radixNode) insertWord(path, fullPath string) {
	rn.path, rn.fullPath = path, fullPath
	rn.passCnt = 1
	rn.end = true
}

// 查看一个单词在 radix 当中是否存在
func (r *Radix) Search(word string) bool {
	node := r.root.search(word)
	return node != nil && node.fullPath == word && node.end
}

// 查找单词流程
func (rn *radixNode) search(word string) *radixNode {
walk:
	for {
		prefix := rn.path
		// word 长于 path

		if len(word) > len(prefix) {
			// 没匹配上，直接返回 nil  word 以节点 path 作为前缀
			if word[:len(prefix)] != prefix {
				return nil
			}
			// word 扣除公共前缀后的剩余部分
			word = word[len(prefix):]
			c := word[0]
			for i := 0; i < len(rn.indices); i++ {
				// 后继节点还有公共前缀，继续匹配 递归开启后续流程
				if c == rn.indices[i] {
					rn = rn.children[i]
					continue walk
				}
			}
			// word 还有剩余部分，但是 prefix 不存在后继节点和 word 剩余部分有公共前缀了
			// 必然不存在
			return nil
		}

		// 和当前节点精准匹配上了
		if word == prefix {
			return rn
		}

		// 走到这里意味着 len(word) <= len(prefix) && word != prefix
		return rn
	}
}

// 前缀匹配流程
func (r *Radix) StartWith(prefix string) bool {
	node := r.root.search(prefix) // 检索出可能包含 prefix 为前缀的节点 node
	// 对应节点存在，并且其全路径 fullPath 确实以 prefix 为前缀，则前缀匹配成功
	return node != nil && strings.HasPrefix(node.fullPath, prefix)
}

// 前缀统计流程
func (r *Radix) PassCnt(prefix string) int {
	node := r.root.search(prefix)
	if node == nil || !strings.HasPrefix(node.fullPath, prefix) {
		return 0
	}
	return node.passCnt // 返回该节点 passCnt 计数器的值
}

// 删除一个单词的流程
func (r *Radix) Erase(word string) bool {
	if !r.Search(word) {
		return false // 判断拟删除单词是否存在，如果不存在直接 return
	}

	// root 直接精准命中了 需要对根节点的所有子节点进行路径 path 调整，同时需要对 radix tree 的根节点指针进行调整
	if r.root.fullPath == word {
		// 如果一个孩子都没有
		if len(r.root.indices) == 0 {
			r.root.path = ""
			r.root.fullPath = ""
			r.root.end = false
			r.root.passCnt = 0
			return true
		}

		// 如果只有一个孩子
		if len(r.root.indices) == 1 {
			r.root.children[0].path = r.root.path + r.root.children[0].path
			r.root = r.root.children[0]
			return true
		}

		// 如果有多个孩子
		for i := 0; i < len(r.root.indices); i++ {
			r.root.children[i].path = r.root.path + r.root.children[0].path
		}

		newRoot := radixNode{
			indices:  r.root.indices,
			children: r.root.children,
			passCnt:  r.root.passCnt - 1,
		}
		r.root = &newRoot
		return true
	}

	// 确定 word 存在的情况下
	move := r.root // 从根节点出发
	// root 单独作为一个分支处理
	// 其他情况下，需要对孩子进行处理
walk:
	for {
		move.passCnt-- // 沿路将途径到的子节点的 passCnt 计数器数值减 1
		prefix := move.path
		word = word[len(prefix):]
		c := word[0]
		// 发现某个子节点的 passCnt 被减为 0，则直接删除该节点. 删除某个子节点后，需要判断，当前节点是否满足和下一个子节点进行压缩合并的条件，如果的是话，需要执行合并操作
		for i := 0; i < len(move.indices); i++ {
			if move.indices[i] != c {
				continue
			}

			// 精准命中但是他仍有后继节点
			if move.children[i].path == word && move.children[i].passCnt > 1 {
				move.children[i].end = false
				move.children[i].passCnt--
				return true
			}

			// 找到对应的 child 了
			// 如果说后继节点的 passCnt = 1，直接干掉
			if move.children[i].passCnt > 1 {
				move = move.children[i]
				continue walk
			}

			move.children = append(move.children[:i], move.children[i+1:]...)
			move.indices = move.indices[:i] + move.indices[i+1:]
			// 如果干掉一个孩子后，发现只有一个孩子了，并且自身 end 为 false 则需要进行合并
			if !move.end && len(move.indices) == 1 {
				// 合并自己与唯一的孩子
				move.path += move.children[0].path
				move.fullPath = move.children[0].fullPath
				move.end = move.children[0].end
				move.indices = move.children[0].indices
				move.children = move.children[0].children
			}

			return true
		}
	}
}

func TestRadixTree(t *testing.T) {
	// 测试压缩前缀树
	trie := NewRadix()
	trie.Insert("/search/v1")
	trie.Insert("/search/v2")
	trie.Insert("/apple")
	trie.Insert("/app")
	fmt.Println(trie.Search("apple"))  // false
	fmt.Println(trie.Search("app"))    // false
	fmt.Println(trie.StartWith("app")) // false
	fmt.Println(trie.Search("/app"))   // true

	trie.Insert("boss")
	fmt.Println(trie.StartWith("bo")) // true
}
