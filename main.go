package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"panda/webframework"
	"regexp"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const LUPAFILE_NAME = "Lupafile"

type NodeState int

const (
	UNVISITED NodeState = iota
	SELECTED
	EXECUTING
	READY
)

type LupaNode struct {
	Target       *LupaTarget
	Mutex        *sync.Mutex
	Dependencies []*LupaNode
	Children     []*LupaNode
	State        NodeState
}

type LupaTarget struct {
	Name     string
	FileDeps []string
	LupaDeps []string
	Script   string
}

var validTargetRe *regexp.Regexp = regexp.MustCompile(`^[A-Za-z0-9\._]+$`)

func isValidTarget(str string) bool {
	return validTargetRe.MatchString(str)
}

type Edge struct {
	target string
	dep    string
}

var allReady chan int

func nodifyTargets(targets []*LupaTarget) (map[string]*LupaNode, error) {
	nodeMap := make(map[string]*LupaNode, 0)

	edges := make([]Edge, 0)

	for _, target := range targets {
		nodeMap[target.Name] = &LupaNode{Target: target, Mutex: &sync.Mutex{}, Dependencies: make([]*LupaNode, 0), Children: make([]*LupaNode, 0), State: UNVISITED}
		for _, dep := range target.LupaDeps {
			edges = append(edges, Edge{target: target.Name, dep: dep})
		}
	}

	for _, edge := range edges {
		target := nodeMap[edge.target]
		dep := nodeMap[edge.dep]
		target.Dependencies = append(target.Dependencies, dep)
		dep.Children = append(dep.Children, target)
	}

	return nodeMap, nil
}

func parseLupafile(content string) ([]*LupaTarget, error) {
	var targets []*LupaTarget
	lines := strings.Split(content, "\n")
	var currentTarget *LupaTarget
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			// bash script
			line = strings.TrimSpace(line) // remove the whitespace
			currentTarget.Script += line + "\n"
		} else {
			// new target declaration
			lineSplit := strings.Split(line, ":")
			if currentTarget != nil {
				targets = append(targets, currentTarget)
			}
			currentTarget = new(LupaTarget)
			currentTarget.Name = lineSplit[0]
			deps := lineSplit[1]
			for _, dep := range strings.Split(deps, " ") {
				if dep == "" {
					continue
				}
				if isValidTarget(dep) {
					currentTarget.LupaDeps = append(currentTarget.LupaDeps, dep)
				} else {
					currentTarget.FileDeps = append(currentTarget.FileDeps, dep)
				}
			}
		}
	}
	if currentTarget != nil {
		targets = append(targets, currentTarget)
	}
	return targets, nil
}

func main() {
	logrus.SetFormatter(&webframework.ErrorFormatter{})
	b, err := ioutil.ReadFile(LUPAFILE_NAME)
	if err != nil {
		switch v := err.(type) {
		case *os.PathError:
			fmt.Println("Lupafile not found!")
			os.Exit(1)
		default:
			panic(v)
		}
	}
	targets, err := parseLupafile(string(b))
	if err != nil {
		panic(err)
	}

	nodes, err := nodifyTargets(targets)
	if err != nil {
		panic(err)
	}

	concurrency := true
	userTarget := "all"
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			// parse arg
			if arg == "-s" {
				concurrency = false
			}
		} else {
			userTarget = arg
		}
	}
	_ = concurrency
	firstNode, ok := nodes[userTarget]
	if !ok {
		fmt.Println("could not find target: " + userTarget)
		os.Exit(1)
	}
	traverse(firstNode)
	/*
		for _, node := range nodes {
			fmt.Println("TARGET", node.Target.Name)
			for _, dep := range node.Dependencies {
				fmt.Println("DEP", dep.Target.Name)
			}
			for _, child := range node.Children {
				fmt.Println("CHILD", child.Target.Name)
			}
			fmt.Println()
		}*/
	/*for _, target := range targets {
		fmt.Println("TARGET NAME", target.Name)
		fmt.Println("FILE DEPS", target.FileDeps)
		fmt.Println("LUPA DEPS", target.LupaDeps)
		fmt.Printf("SCRIPT:\n%s\n", target.Script)
	}*/
}

/*

TODO:
add waitgroup
add timestamp check

*/
