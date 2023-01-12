package main

import (
	"fmt"
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
	handler, ok := h.root.findRouter(c.r.URL.Path)
	if !ok {
		_ = c.BadRequestJson(commonResponse{BizCode: 1, Msg: "路由不存在"})
		return
	}
	handler(c)
}

func (h *HandlerBaseOnTree) Route(method, pattern string, handle handleFunc) {
	paths := strings.Split(strings.Trim(pattern, "/"), "/")

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
	var wildcard *node

	for _, child := range n.children {
		if child.path == path && child.path != "*" {
			return child, true
		}

		if child.path == "*" {
			wildcard = child
		}
	}
	return wildcard, wildcard != nil
}

func (n *node) findRouter(pattern string) (handleFunc, bool) {
	paths := strings.Split(strings.Trim(pattern, "/"), "/")

	cur := n
	for _, path := range paths {
		matchChild, ok := cur.findMatchChild(path)
		if !ok {
			return nil, false
		}
		cur = matchChild
	}

	if cur.handler == nil {
		return nil, false
	}
	return cur.handler, true
}
