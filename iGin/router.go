package iGin

import (
	"log"
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
			log.Printf("newNode edge:<%s>", path)
			nextNode := newNode()
			nextNode.par = now
			nextNode.edge = path
			now.son = append(now.son, nextNode)
			now = nextNode
		}
	}
	return true, now
}

//用于管理url在Trie中的遍历以及动态路由解析
type urlHelper struct {
	index     int
	pathSlice []string
	url       string
	isFinish  bool
}

func newUrlHelper(url string) *urlHelper {
	//log.Printf("buildUrlHelper url:%s", url)
	return &urlHelper{
		index:     -1,
		pathSlice: strings.Split(strings.Trim(url, "/"), "/"),
		url:       url,
		isFinish:  false,
	}
}

func (u *urlHelper) advance() (int, string, bool) {
	u.index++
	if u.index == len(u.pathSlice) {
		u.setFinish(true)
	}
	if u.isFinish || u.index < 0 || u.index >= len(u.pathSlice) {
		return -1, "", false
	}
	return u.index, u.pathSlice[u.index], true
}

func (u *urlHelper) setFinish(value bool) {
	u.isFinish = value
}

func (u *urlHelper) nowPath() string {
	if u.index >= len(u.pathSlice) || u.index < 0 {
		return ""
	}
	return u.pathSlice[u.index]
}

func (u *urlHelper) match(dst string) (map[string]string, bool) {
	if u.nowPath() == dst {
		return nil, true
	}
	if len(dst) > 1 {
		switch dst[0] {
		case ':':
			params := make(map[string]string)
			params[dst[1:]] = u.nowPath()
			return params, true
		case '*':
			params := make(map[string]string)
			params[dst[1:]] = strings.Join(u.pathSlice[u.index:], "/")
			u.setFinish(true)
			return params, true
		}
	}
	return nil, false
}

//询问url是否存在，并返回处理函数
func (r *RouterManager) Query(url string) (bool, []HandlerFunc, map[string]string) {
	urlHp := newUrlHelper(url)
	if r.Root == nil {
		return false, nil, nil
	}
	handlers := make([]HandlerFunc, 0)
	params := make(map[string]string)
	now := r.Root
	for {
		index, path, ok := urlHp.advance()
		log.Printf("urlHp,advance index:%d,path:%s,ok:%v", index, path, ok)
		if !ok {
			break
		}

		find := false
		for _, son := range now.son {
			log.Printf("son.Edge:<%s>,now.son:%+v", son.edge, now.son)
			if param, ok := urlHp.match(son.edge); ok {
				handlers = append(handlers, son.middleWare...)
				for k, v := range param {
					params[k] = v
				}
				find = true
				now = son
				break
			}
		}
		log.Printf("index:%v,url:%v,find:%v,path:%v", index, url, find, path)
		if !find {
			break
		}
	}
	//如果没有视图函数,使用默认函数
	if now.viewFunc == nil {
		now.viewFunc = DefaultNotFound
	}
	handlers = append(handlers, now.viewFunc)
	return true, handlers, params
}
