package gee

import (
	"fmt"
	"strings"
)

// node router 前缀树
// 与普通的树不同, 为了实现动态路由匹配, 加上了 isWild 这个参数
type node struct {
	// pattern 待匹配路由, 例如 /p/:lang
	pattern string

	// part 路由中的一部分, 例如 /:lang
	part string

	// children 子节点, 例如 [doc, tutorial, intro]
	children []*node

	// isWild 是否精确匹配, part 含有 * / : 时为 true
	isWild bool
}

// String 返回 node 信息
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// matchChild 第一个匹配成功的节点, 用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 精确匹配或模糊匹配
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

// matchChildren 所有匹配成功的节点, 用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0, len(part))

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

// insert 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果已经到 parts 的最后一层数据, 则表明插入结束, 将 pattern 放入节点的 pattern 中
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 获取第 height 层 part, 并查询对应 part 是否有子节点
	part := parts[height]
	child := n.matchChild(part)

	// 如果没有子节点, 则构建一个新的子节点
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search 查询节点
func (n *node) search(parts []string, height int) *node {
	// 已经找到 parts 的最后一层数据或 当前 part 前缀为 *(全匹配)
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果当前树节点 pattern == "", 表明该节点并非最后一个节点, 匹配失败
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 获取第 height 层 part, 并查询对应 part 是否有子节点
	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)

		if result != nil {
			return result
		}
	}

	return nil
}

// travel 返回 node 树中的所有 pattern
func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
