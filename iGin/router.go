package iGin

import (
	"fmt"
	"strings"
)

//Trie树上的节点
type node struct {
	edge       string
	par        *node
	son        []*node
	viewFunc   HandlerFunc   //视图函数
	middleWare []HandlerFunc //中间件
}

func newNode() *node {
	return &node{edge: "", par: nil, son: make([]*node, 0)}
}

//管理插入和查询 url到function的映射
type RouterManager struct {
	Root *node
}

func (r *RouterManager) CheckUrlValid(url string) (bool, error) {
	return true, nil
}

func (r *RouterManager) InsertViewFunc(url string, handler HandlerFunc) (bool, error) {
	if ok, now := r.getTrieNode(url); ok {
		now.viewFunc = handler
		return true, nil
	}
	return false, nil
}

func (r *RouterManager) InsertMiddleWare(url string, handlers []HandlerFunc) (bool, error) {
	if ok, now := r.getTrieNode(url); ok {
		now.middleWare = append(now.middleWare, handlers...)
		return true, nil
	}
	return false, nil
}

//返回url在Trie上对应的node节点
func (r *RouterManager) getTrieNode(url string) (bool, *node) {
	if valid, _ := r.CheckUrlValid(url); !valid {
		return false, nil
	}
	pathStr := strings.Split(strings.Trim(url, "/"), "/")
	if r.Root == nil {
		r.Root = newNode()
	}
	now := r.Root
	for _, path := range pathStr {
		fmt.Println(path)
		hasSon := false
		for _, son := range now.son {
			if son.edge == path {
				hasSon = true
				now = son
				break
			}
		}
		//如果没有子节点,那么新建一个
		if !hasSon {
			nextNode := newNode()
			nextNode.par = now
			nextNode.edge = path
			now.son = append(now.son, nextNode)
			now = nextNode
		}
	}
	return true, now
}

//询问url是否存在，并返回处理函数
func (r *RouterManager) Query(url string) (bool, []HandlerFunc) {
	pathStr := strings.Split(strings.Trim(url, "/"), "/")
	if r.Root == nil {
		return false, nil
	}
	handlers := make([]HandlerFunc, 0)
	now := r.Root
	for _, path := range pathStr {
		find := false
		for _, son := range now.son {
			if son.edge == path {
				find = true
				handlers = append(handlers, son.middleWare...)
				now = son
				break
			}
		}
		if !find {
			return false, nil
		}
	}
	handlers = append(handlers, now.viewFunc)
	return true, handlers
}
