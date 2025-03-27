package misc

import (
	"satellite/resources"
	"strings"
)

var acAutomaton *ACAutomaton

type TrieNode struct {
	children map[rune]*TrieNode
	fail     *TrieNode
	output   []string
}

type ACAutomaton struct {
	root *TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		fail:     nil,
		output:   []string{},
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
	for _, existing := range node.output {
		if existing == pattern {
			return
		}
	}
	node.output = append(node.output, pattern)
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
			child.output = append(child.output, child.fail.output...)
			queue = append(queue, child)
		}
	}
}

func (ac *ACAutomaton) Search(text string) map[string][]int {
	result := make(map[string][]int)
	node := ac.root
	for i, ch := range text {
		for node != ac.root && node.children[ch] == nil {
			node = node.fail
		}
		if next, exists := node.children[ch]; exists {
			node = next
		} else {
			node = ac.root
		}
		for _, pattern := range node.output {
			result[pattern] = append(result[pattern], i-len(pattern)+1)
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
