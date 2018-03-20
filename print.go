package main

import (
	"fmt"
)

func printNode(node *LupaNode) {
	deps := ""
	for _, dep := range node.Dependencies {
		deps += dep.Target.Name + ", "
	}
	children := ""
	for _, child := range node.Children {
		children += child.Target.Name + ", "
	}
	fmt.Printf("(%s | DEPS %s | CHILDREN %s)\n", node.Target.Name, deps, children)
}
