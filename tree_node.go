package main

import "strings"

/**
节点匹配类型 优先级按照 大->小
*/
const (
	// 根
	nodeTypeRoot = iota

	// *
	nodeTypeAny

	// 路径参数
	nodeTypeParam

	// 全量匹配
	nodeTypeStatic
)

const _any = "*"

type matchFunc func(path string, ctx *Context) bool

type node struct {
	children []*node
	pattern  string

	matchFunc matchFunc
	handler   handleFunc

	nodeType int32
}

func newNode(path string) *node {
	if path == _any {
		return newAnyNode()
	}
	if strings.HasPrefix(path, ":") {
		return newParamNode(path)
	}
	return newStaticNode(path)
}

func newRootNode(method string) *node {
	return &node{
		children: make([]*node, 0),
		pattern:  method,
		matchFunc: func(p string, ctx *Context) bool {
			return true
		},
		nodeType: nodeTypeRoot,
	}
}

func newAnyNode() *node {
	return &node{
		pattern: _any,
		matchFunc: func(p string, ctx *Context) bool {
			return true
		},
		nodeType: nodeTypeAny,
	}
}

func newParamNode(path string) *node {
	name := path[1:]
	return &node{
		children: make([]*node, 0),
		pattern:  path,
		matchFunc: func(p string, ctx *Context) bool {
			if ctx != nil {
				ctx.PathParams[name] = p
			}
			return p != _any
		},
		nodeType: nodeTypeParam,
	}
}

func newStaticNode(path string) *node {
	return &node{
		children: make([]*node, 0),
		pattern:  path,
		matchFunc: func(p string, ctx *Context) bool {
			return path == p && path != _any
		},
		nodeType: nodeTypeStatic,
	}
}
