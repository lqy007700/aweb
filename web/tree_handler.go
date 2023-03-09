package web

import (
	"net/http"
	"strings"
)

var supperMethods = [4]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodDelete,
	http.MethodPut,
}

type HandlerOnTree struct {
	paths map[string]*node
}

func NewHandlerOnTree() Handler {
	paths := make(map[string]*node, len(supperMethods))
	for _, method := range supperMethods {
		paths[method] = newRootNode()
	}
	return &HandlerOnTree{paths: paths}
}

type node struct {
	children []*node
	handler  HandlerFunc
	pattern  string
	midls    []Middleware
}

func newRootNode() *node {
	return &node{
		children: make([]*node, 0),
	}
}
func newNode(path string) *node {
	return &node{
		children: make([]*node, 0),
		pattern:  path,
	}
}

func (h *HandlerOnTree) ServeHTTP(ctx *Context) {
	n, ok := h.findRouter(ctx)
	if !ok {
		ctx.BadJson(commonResp{
			Code: -1,
			Msg:  ErrorRouterNotFound.Error(),
			Data: nil,
		})
		return
	}

	// middleware
	// 最后一个应该是执行用户代码
	var root HandlerFunc = func(ctx *Context) {
		if !ok || n == nil || n.handler == nil {
			ctx.BadJson(commonResp{
				Code: -1,
				Msg:  ErrorRouterNotFound.Error(),
				Data: nil,
			})
			return
		}
		n.handler(ctx)
	}

	for i := len(n.midls) - 1; i >= 0; i-- {
		root = n.midls[i](root)
	}
	root(ctx)
}

func (h *HandlerOnTree) Route(method, pattern string, handler HandlerFunc, ms ...Middleware) {
	cur, ok := h.paths[method]
	if !ok {
		panic(ErrorInvalidMethod)
	}
	paths := strings.Split(strings.Trim(pattern, "/"), "/")
	for i, path := range paths {
		child, ok := h.findMatchChild(cur, path)
		if ok {
			cur = child
		} else {
			h.createSubTree(cur, paths[i:], handler, ms...)
			return
		}
	}
}

// 查找路由
func (h *HandlerOnTree) findRouter(ctx *Context) (*node, bool) {
	path := ctx.R.URL.Path
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur, ok := h.paths[ctx.R.Method]
	if !ok {
		return nil, false
	}

	for _, p := range paths {
		child, ok := h.findMatchChild(cur, p)
		if !ok {
			return nil, false
		}
		cur = child
	}
	return cur, true
}

func (h *HandlerOnTree) findMatchChild(cur *node, p string) (*node, bool) {
	for _, child := range cur.children {
		if child.pattern == p {
			return child, true
		}
	}
	return nil, false
}

// 创建节点
func (h *HandlerOnTree) createSubTree(root *node, patterns []string, handler HandlerFunc, ms ...Middleware) {
	cur := root
	for _, path := range patterns {
		n := newNode(path)
		cur.children = append(cur.children, n)
		cur = n
	}
	cur.handler = handler
	cur.midls = ms
}
