package pre_tree_test

import (
	"fmt"
	"testing"
)

type trieNode struct {
	nexts   [26]*trieNode
	passCnt int // 保证在后续处理单词插入的 Insert 流程以及单词删除的 Erase 流程中，对每个节点维护好一个 passCnt 计数器，用于记录通过该节点的单词数量.
	end     bool
}

type Trie struct {
	root *trieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &trieNode{},
	}
}

// Search 查找是否含有整个单词
func (t *Trie) Search(word string) bool {
	// 查找目标节点，使得从根节点开始抵达目标节点沿路字符形成的字符串恰好等于 word
	node := t.search(word)
	// tag: 如果节点存在，并且节点的 end 标识为 true，代表 word 存在
	return node != nil && node.end
}

// 查找字符串的节点的底层 后续跟进end标识判断
func (t *Trie) search(target string) *trieNode {
	// 移动指针从根节点出发
	move := t.root
	// 依次遍历 target 中的每个字符
	for _, ch := range target {
		// 倘若 nexts 中不存在对应于这个字符的节点，说明该单词没插入过，返回 nil
		if move.nexts[ch-'a'] == nil {
			return nil
		}
		// 指针向着子节点移动
		move = move.nexts[ch-'a']
	}

	// 来到末尾，说明已经完全匹配好单词，直接返回这个节点
	// 需要注意，找到目标节点不一定代表单词存在，因为该节点的 end 标识未必为 true
	// 比如我们之前往 trie 中插入了 apple 这个单词，但是查找 app 这个单词时，预期的返回结果应该是不存在，此时就需要使用到 end 标识 进行区分
	return move
}

// StartWith 查找是否含有单词前缀
func (t *Trie) StartWith(prefix string) bool {
	// StartWith 无需对节点的 end 标识进行判断
	return t.search(prefix) != nil
}

// PassCnt 给定一个 prefix，要求统计出以 prefix 作为前缀的单词数量
func (t *Trie) PassCnt(prefix string) int {
	node := t.search(prefix)
	if node == nil {
		return 0
	}

	return node.passCnt
}

// Erase 从前缀树 trie 中插入某个单词
func (t *Trie) Insert(word string) {
	if t.Search(word) {
		return
	}
	// 从根节点开始
	move := t.root
	// 依次遍历 word 的每个字符，每轮判断当前节点的子节点列表 nexts 中，对应于字符的子节点是否已存在了
	for _, ch := range word {
		if move.nexts[ch-'a'] == nil {
			move.nexts[ch-'a'] = &trieNode{}
		}
		move.nexts[ch-'a'].passCnt++ // 存在新单词的插入，我们需要对这个子节点的 passCnt 计数器累加 1
		move = move.nexts[ch-'a']
	}
	// 此时 move 所在位置一定对应的是单词结尾的字符. 我们需要把 move 指向节点的 end 标识置为 true，代表存在单词以此节点作为结尾
	move.end = true
}

// Erase 从前缀树 trie 中删除某个单词
func (t *Trie) Erase(word string) bool {
	if !t.Search(word) {
		return false
	}

	move := t.root
	// 依次遍历 word 中的每个字符，每次从当前节点子节点列表 nexts 中找到对应于字符的子节点
	for _, ch := range word {
		move.nexts[ch-'a'].passCnt--
		// 倘若发现子节点的 passCnt 被减为 0，则直接舍弃这个子节点，结束流程
		if move.nexts[ch-'a'].passCnt == 0 {
			move.nexts[ch-'a'] = nil
			return true
		}
		move = move.nexts[ch-'a']
	}
	// 遍历来到单词末尾位置，则需要把对应节点的 end 标识置为 false
	move.end = false
	return true
}

func TestPreTree(t *testing.T) {
	// 测试前缀树
	trie := NewTrie()
	trie.Insert("apple")
	fmt.Println(trie.Search("apple")) // 返回 True
	fmt.Println(trie.Search("app"))   // 返回 False

	trie.Insert("app")
	fmt.Println(trie.Search("app")) // 返回 True

	trie.Insert("boss")
	fmt.Println(trie.StartWith("bo")) // 返回 True
	fmt.Println(trie.StartWith("ab")) // 返回 false
}
