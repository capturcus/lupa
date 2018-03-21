package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"panda/webframework"
	"strings"
	"sync"
	"unicode/utf8"
)

func spacePad(numSpaces int, targetName, message string) string {
	return targetName + strings.Repeat(" ", numSpaces-utf8.RuneCountInString(targetName)) + ": " + message
}

func printPipe(targetName string, pipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(spacePad(maxTargetLength, targetName, m))
	}
}

func executeScript(targetName string, script string) error {
	cmd := exec.Command("bash", "-s")
	cmd.Stdin = strings.NewReader(script)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return webframework.ReturnError(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go printPipe(targetName, stdout, wg)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return webframework.ReturnError(err)
	}
	wg.Add(1)
	go printPipe(targetName, stderr, wg)
	err = cmd.Start()
	if err != nil {
		return webframework.ReturnError(err)
	}
	wg.Wait() // for the pipes
	err = stderr.Close()
	if err != nil {
		panic(err)
	}
	err = stdout.Close()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		return webframework.ReturnError(err)
	}
	return nil
}

func executeNode(node *LupaNode, wg *sync.WaitGroup) {
	defer wg.Done()
	// fmt.Printf("executing %s\n", node.Target.Name)
	err := executeScript(node.Target.Name, node.Target.Script)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	node.Mutex.Lock()
	node.State = READY
	node.Mutex.Unlock()
	for _, child := range node.Children {
		checkAndExecute(child, wg)
	}
}

func checkAndExecute(node *LupaNode, wg *sync.WaitGroup) {
	// fmt.Printf("checking %s\n", node.Target.Name)
	if node.State != SELECTED {
		// fmt.Printf("node %s not selected\n", node.Target.Name)
		return
	}
	node.Mutex.Lock()
	// fmt.Printf("checking lock %s\n", node.Target.Name)
	for _, dep := range node.Dependencies {
		if dep.State != READY {
			node.Mutex.Unlock()
			// fmt.Printf("node %s not ready\n", node.Target.Name)
			return
		}
	}
	node.State = EXECUTING
	wg.Add(1)
	go executeNode(node, wg)
	node.Mutex.Unlock()
}

func traverse(node *LupaNode, wg *sync.WaitGroup) {
	if node.State != UNVISITED {
		// we were here already
		// fmt.Printf("node %s already visited\n", node.Target.Name)
		return
	}
	// fmt.Println("traversing\n", node.Target.Name)
	node.Mutex.Lock()
	node.State = SELECTED
	if len(node.Dependencies) == 0 {
		node.State = EXECUTING
		wg.Add(1)
		// fmt.Printf("node %s is leaf\n", node.Target.Name)
		go executeNode(node, wg)
	}
	node.Mutex.Unlock()
	for _, dep := range node.Dependencies {
		traverse(dep, wg)
	}
}
