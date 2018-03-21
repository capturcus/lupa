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

func printTargets(targets []*LupaTarget) {
	for _, target := range targets {
		fmt.Println("TARGET NAME", target.Name)
		fmt.Println("FILE DEPS", target.FileDeps)
		fmt.Println("LUPA DEPS", target.LupaDeps)
		fmt.Printf("SCRIPT:\n%s\n", target.Script)
	}
}
