package main

import (
	"strings"
)

type HandlerBaseOnTree struct {
	root *node
}

func NewHandlerBaseOnTree() *HandlerBaseOnTree {
	return &HandlerBaseOnTree{root: newNode("/")}
}

type node struct {
	path     string
	children []*node

	handler handleFunc
}

func newNode(path string) *node {
	return &node{path: path, children: make([]*node, 0)}
}

func (h *HandlerBaseOnTree) ServeHTTP(c *Context) {
	node := h.root.matchRoute(c.r.URL.Path)
	if node != nil && node.handler != nil {
		node.handler(c)
	} else {
		err := c.BadRequestJson(commonResponse{
			BizCode: 1,
			Msg:     "路由不存在",
		})
		if err != nil {
			return
		}
	}
}

func (h *HandlerBaseOnTree) Route(method, pattern string, handle handleFunc) {
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	cur := h.root
	for idx, path := range paths {
		matchChild, ok := cur.findMatchChild(path)
		if ok {
			cur = matchChild
		} else {
			cur.createSubTree(paths[idx:], handle)
			return
		}
	}
}

// 创建节点 绑定handle
func (n *node) createSubTree(path []string, handle handleFunc) {
	cur := n
	for _, s := range path {
		nn := newNode(s)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handle
}

func (n *node) findMatchChild(path string) (*node, bool) {
	for _, child := range n.children {
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}

func (n *node) matchRoute(pattern string) *node {
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	cur := n
	for _, path := range paths {
		matchChild, ok := cur.findMatchChild(path)
		if ok {
			cur = matchChild
		} else {
			cur = nil
		}
	}
	return cur
}
