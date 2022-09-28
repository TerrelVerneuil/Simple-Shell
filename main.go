package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Terminal contains
type Terminal struct {
}

// Execute executes a given command
func (t *Terminal) Execute(command string) {

}

// This is the main function of the application.
// User input should be continuously read and checked for commands
// for all the defined operations.
// See https://golang.org/pkg/bufio/#Reader and especially the ReadLine
// function.
func main() {
	wdir, _ := os.Getwd()
	fmt.Print(wdir + "> ")
	reader := bufio.NewReader(os.Stdin)
	var his []string
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		switch {
		case strings.HasPrefix(text, "cd"):
			cd := os.Chdir(strings.SplitAfter(text, " ")[1])
			wdir, _ = os.Getwd()
			if cd != nil {
				fmt.Print()
			}
		}
		cmd := exec.Command(string(text))
		cmd.Dir = wdir
		out, _ := cmd.Output()
		err := cmd.Run()
		his = append(his, text)
		switch {
		case strings.HasPrefix(text, "history"):
			for i := 0; i < len(his); i++ {
				fmt.Println(string(rune(i)) + " " + his[i])
			}
		case strings.HasPrefix(text, "ls "):
			files, _ := ioutil.ReadDir(strings.SplitAfter(text, " ")[1])
			for _, f := range files {
				fmt.Println(f.Name())
			}
			if err != nil {
				fmt.Println(err)
			}
		case strings.HasPrefix(text, "cat"):
			content, _ := os.ReadFile(strings.SplitAfter(text, " ")[1])
			fmt.Println(string(content))
			if err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(text, "exit"):
			os.Exit(0)
			if err != nil {
				fmt.Println(err)
			}

		case strings.HasPrefix(text, "mkdir"):
			mkdir := os.Mkdir(strings.SplitAfter(text, " ")[1], 0750)
			if mkdir != nil {
				fmt.Println()
			}
		case strings.HasPrefix(text, "rm"):
			os.Remove(strings.SplitAfter(text, " ")[1])
			if err != nil {
				fmt.Println(err)
			}
		case strings.HasPrefix(text, "create"):
			create, _ := os.Create(strings.SplitAfter(text, " ")[1])
			if create != nil {
				fmt.Println()
			}

		case strings.HasPrefix(text, "rm -r"):
			os.RemoveAll(strings.SplitAfter(text, " ")[1])
			if err != nil {
				fmt.Println(err)
			}

		}

		if err != nil {
			fmt.Println(err)
		}

		fmt.Print(string(out))
		fmt.Print(wdir + "> ")
	}
}
