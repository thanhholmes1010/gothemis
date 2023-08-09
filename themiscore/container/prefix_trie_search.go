package container

import (
	"fmt"
	"strings"
)

// PrefixNode
type PrefixNode struct {
	key            string
	childs         []*PrefixNode
	allowPassedAll bool
	fullKey        string
	data           any
}

func (n *PrefixNode) matchChild(key string) *PrefixNode {
	for _, child := range n.childs {
		if child.key == key {
			return child
		}
	}
	return nil
}

func newPrefixNode(key string) *PrefixNode {
	return &PrefixNode{
		key: key,
	}
}

type PrefixTrie struct {
	root            *PrefixNode
	defaultWildCard string
	defaultSep      string
}

func NewPrefixTrie(defaultSep string, defaultWildCard string) *PrefixTrie {
	return &PrefixTrie{
		defaultSep:      defaultSep,
		defaultWildCard: defaultWildCard,
		root: &PrefixNode{
			key: "",
		},
	}
}

func (p *PrefixTrie) Insert(fullKey string, data any) {
	n := p.root
	keys := strings.Split(fullKey, p.defaultSep)
	fmt.Println(keys, len(keys))
	for _, key := range keys {
		child := n.matchChild(key)
		if child == nil {
			child = newPrefixNode(key)
			if strings.HasPrefix(key, p.defaultWildCard) {
				child.allowPassedAll = true
			}
			n.childs = append(n.childs, child)
		}
		n = child
	}
	// n here is depth child
	n.fullKey = fullKey
	n.data = data
}

type color uint8

const (
	white color = iota // not visited
	black              // backtracked
)

type PairPrefixNodeData struct {
	node  *PrefixNode
	color color
}

func (p *PrefixTrie) SearchAllPossibles(fullKey string, funcDraw func(treeRecommendInfo *string, node *PrefixNode,
	level int, addData *string)) any {

	keys := strings.Split(fullKey, p.defaultSep)
	l := len(keys)
	// implement backtracking
	stack := NewQueue()
	stack.Push(false, p.root)
	index := 0
	level := -1
	recommendQueue := NewQueue()
	for !stack.Empty() {

		size := stack.Size()
		//findLevel := false
		for size > 0 {
			n := stack.Pop().(*PrefixNode)
			if index < l {
				for _, child := range n.childs {
					if child.key == keys[index] || child.allowPassedAll {
						stack.Push(false, child)
						if index == l-1 && child.data != nil {
							return child.data
						}
					}
				}
			}
			if n.key != "" {
				recommendQueue.Push(false, &drawData{
					n:     n,
					level: level,
				})
			}
			size--
		}
		level++
		index++
	}
	p.Recommend(&keys, recommendQueue, funcDraw)
	return nil
}

type drawData struct {
	n     *PrefixNode
	level int
}

func (p *PrefixTrie) Recommend(fullKey *[]string, recommendQueue IteratorContainer, funcDraw func(treeRecommendInfo *string, node *PrefixNode, level int, addData *string)) {
	treeRecomendInfo := ""
	level := -1
	matchEnd := false
	//for !recommendQueue.Empty() {
	//	fmt.Println("key: ", recommendQueue.Pop().(*PrefixNode).key)
	//}
	index := 0
	l := len(*fullKey)
	for !recommendQueue.Empty() {
		size := recommendQueue.Size()
		for size > 0 {

			node := recommendQueue.Pop().(*drawData)
			fmt.Println("node key: ", node.n.key, " level: ", node.level)
			if !matchEnd {
				key := ""
				if index < l {
					if node.n.allowPassedAll {
						key = "[param: " + (*fullKey)[index] + "]"
					}
					if index == l-1 {
						key = "[last match key: " + (*fullKey)[index] + "]" + ", do you mean ?"
					}
					if key != "" {
						funcDraw(&treeRecomendInfo, node.n, node.level, &key)
					} else {
						funcDraw(&treeRecomendInfo, node.n, node.level, nil)
					}
				}
				level++
				index++
				if index == l {
					for _, child := range node.n.childs {
						recommendQueue.Push(false, &drawData{
							child,
							node.level + 1,
						})
					}
				}
			} else {
				fmt.Println("first key in here: ", node.n.key)
				funcDraw(&treeRecomendInfo, node.n, node.level, nil)
				for _, child := range node.n.childs {
					fmt.Println("parent: ", node.n.key, " child: ", child.key)
					recommendQueue.Push(true, &drawData{
						n:     child,
						level: node.level + 1,
					})
				}
			}
			size--
		}
		if !matchEnd {
			matchEnd = true
		}
		if matchEnd {
			level++
		}
	}
	fmt.Println(treeRecomendInfo)
}
func (p *PrefixTrie) DrawLevelWithTabSpace(treeRecommendInfo *string, node *PrefixNode, level int, addData *string) {
	for i := 0; i < level; i++ {
		*treeRecommendInfo += "    "
	}
	*treeRecommendInfo += fmt.Sprintf("-> %v", node.key)
	if addData != nil {
		*treeRecommendInfo += *addData
	}
	*treeRecommendInfo += "\n"
}
