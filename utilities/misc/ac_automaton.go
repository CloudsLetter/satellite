package misc

import (
	"satellite/resources"
	"strings"
)

var acAutomaton *ACAutomaton

type TrieNode struct {
	children map[rune]*TrieNode
	fail     *TrieNode
	output   map[string]struct{}
}

type ACAutomaton struct {
	root *TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		fail:     nil,
		output:   make(map[string]struct{}),
	}
}

func NewACAutomaton() *ACAutomaton {
	return &ACAutomaton{root: NewTrieNode()}
}

func (ac *ACAutomaton) Insert(pattern string) {
	node := ac.root
	for _, ch := range pattern {
		if _, exists := node.children[ch]; !exists {
			node.children[ch] = NewTrieNode()
		}
		node = node.children[ch]
	}
	node.output[pattern] = struct{}{}
}

func (ac *ACAutomaton) Build() {
	queue := []*TrieNode{}

	for _, child := range ac.root.children {
		child.fail = ac.root
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for ch, child := range current.children {
			failNode := current.fail
			for failNode != nil {
				if next, exists := failNode.children[ch]; exists {
					child.fail = next
					break
				}
				failNode = failNode.fail
			}
			if failNode == nil {
				child.fail = ac.root
			}
			for pattern := range child.fail.output {
				child.output[pattern] = struct{}{}
			}
			queue = append(queue, child)
		}
	}
}

func (ac *ACAutomaton) Search(text string) map[string][]int {
	result := make(map[string][]int)
	node := ac.root
	runeText := []rune(text) // 处理 Unicode 字符正确索引
	for i, ch := range runeText {
		for node != ac.root && node.children[ch] == nil {
			node = node.fail
		}
		if next, exists := node.children[ch]; exists {
			node = next
		} else {
			node = ac.root
		}
		for pattern := range node.output {
			startIndex := i - len([]rune(pattern)) + 1
			if startIndex >= 0 {
				result[pattern] = append(result[pattern], startIndex)
			}
		}
	}
	return result
}

func InitACAutomaton() {
	ac := NewACAutomaton()
	wordArr := strings.Split(resources.GetSensitiveWord(), "\n")
	for _, word := range wordArr {
		ac.Insert(word)
	}
	ac.Build()
	acAutomaton = ac
	return
}

func GetACAutomaton() *ACAutomaton {
	return acAutomaton
}
