package main

import (
	"errors"
	"net/http"
	"strings"
)

var ErrorInvalidRouterPattern = errors.New("invalid router pattern")
var ErrorInvalidMethod = errors.New("invalid method")
var ErrorRouterNotFound = errors.New("router not found")

var supportMethods = [4]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
}

type HandlerBaseOnTree struct {
	forest map[string]*node
}

func NewHandlerBaseOnTree() *HandlerBaseOnTree {
	forests := make(map[string]*node, len(supportMethods))
	for _, m := range supportMethods {
		forests[m] = newRootNode(m)
	}

	return &HandlerBaseOnTree{
		forest: forests,
	}
}

func (h *HandlerBaseOnTree) ServeHTTP(c *Context) {
	handler, ok := h.root.findRouter(c.r.URL.Path)
	if !ok {
		_ = c.BadRequestJson(commonResponse{BizCode: 1, Msg: "路由不存在"})
		return
	}
	handler(c)
}

func (h *HandlerBaseOnTree) Route(method, pattern string, handle handleFunc) error {
	err := h.validatePattern(pattern)
	if err != nil {
		return err
	}

	paths := strings.Split(strings.Trim(pattern, "/"), "/")

	cur, ok := h.forest[method]
	if !ok {
		return ErrorInvalidMethod
	}

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

/**
检测路由规范
1/ 校验 *，如果存在，必须在最后一个，并且它前面必须是/
*/
func (h *HandlerBaseOnTree) validatePattern(pattern string) error {
	index := strings.Index(pattern, _any)
	if index > 0 {
		if index != len(pattern)-1 {
			return ErrorInvalidRouterPattern
		}

		if pattern[index-1] != '/' {
			return ErrorInvalidRouterPattern
		}
	}
	return nil
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
