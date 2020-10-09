package yp

import (
	"golang.org/x/net/html"
	"strings"
)

type node struct {
	n html.Node
}

func (n *node) getAttr(key string) string {
	for _, attr := range n.n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func (n *node) getFirstChildNode() *node {
	for c := n.n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			return &node{*c}
		}
	}
	return nil
}

func (n *node) getClassesS(class string) []*node {
	var classes []*node

	var r func(nn *node)
	r = func(nn *node) {
		if nn.n.Type == html.ElementNode && nn.getAttr("class") == class {
			classes = append(classes, nn)

		}
		for c := nn.n.FirstChild; c != nil; c = c.NextSibling {
			r(&node{*c})
		}
	}

	r(n)

	return classes
}

func (n *node) getClasses(class string) []*node {
	var classes []*node

	var r func(nn *node)
	r = func(nn *node) {
		if nn.n.Type == html.ElementNode {
			for _, cl := range strings.Split(nn.getAttr("class"), " ") {
				if cl == class {
					classes = append(classes, nn)
					break
				}
			}

		}

		for c := nn.n.FirstChild; c != nil; c = c.NextSibling {
			r(&node{*c})
		}
	}

	r(n)

	return classes
}

func (n *node) getChildrenOfType(t string) []*node {
	var nodes []*node
	for c := n.n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == t {
			nodes = append(nodes, &node{*c})
		}
	}

	return nodes
}
