package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"panda/webframework"
	"strings"
)

func printPipe(targetName string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(targetName+":", m)
	}
}

func executeScript(targetName string, script string) error {
	cmd := exec.Command("bash", "-s")
	cmd.Stdin = strings.NewReader(script)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return webframework.ReturnError(err)
	}
	go printPipe(targetName, stdout)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return webframework.ReturnError(err)
	}
	go printPipe(targetName, stderr)
	err = cmd.Start()
	if err != nil {
		return webframework.ReturnError(err)
	}
	err = cmd.Wait()
	if err != nil {
		return webframework.ReturnError(err)
	}
	return nil
}

func executeNode(node *LupaNode) {
	fmt.Println("executing", node.Target.Name)
	err := executeScript(node.Target.Name, node.Target.Script)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	node.Mutex.Lock()
	node.State = READY
	node.Mutex.Unlock()
	for _, child := range node.Children {
		go checkAndExecute(child)
	}
}

func checkAndExecute(node *LupaNode) {
	if node.State != SELECTED {
		return
	}
	node.Mutex.Lock()
	for _, dep := range node.Dependencies {
		if dep.State != READY {
			node.Mutex.Unlock()
			return
		}
	}
	node.State = EXECUTING
	go executeNode(node)
	node.Mutex.Unlock()
}

func traverse(node *LupaNode) {
	if node.State != UNVISITED {
		// we were here already
		return
	}
	fmt.Println("traversing", node.Target.Name)
	node.Mutex.Lock()
	node.State = SELECTED
	if len(node.Dependencies) == 0 {
		node.State = EXECUTING
		go executeNode(node)
	}
	node.Mutex.Unlock()
	for _, dep := range node.Dependencies {
		go traverse(dep)
	}
}
