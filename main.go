package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"net/http"
    "github.com/gorilla/websocket"
    "log"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    for {
        // Read message as JSON and map it to a Message object
        var msg string
        err := ws.ReadJSON(&msg)
        if err != nil {
            log.Printf("error: %v", err)
            break
        }
        // Here you can handle the message (execute command)
        log.Printf("Received: %s\n", msg)
        // Respond back
        err = ws.WriteJSON("Executed: " + msg)
        if err != nil {
            log.Printf("error: %v", err)
            break
        }
    }
}

func main() {
	go startWebServer()
	terminal()
}
func terminal() {
    reader := bufio.NewReader(os.Stdin)
    var history []string
    currentDir, _ := os.Getwd()

    for {
        fmt.Print(currentDir + "> ")
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)

        history = append(history, text)
        command, args := parseCommand(text)

		switch command {
		case "cd":
			currentDir = changeDirectory(args, currentDir)
		case "ls":
			listFiles(args)
		case "cat":
			catFile(args)
		case "history":
			showHistory(history)
		case "mkdir":
			makeDirectory(args)
		case "rm":
			removeFile(args)
		case "exit":
			os.Exit(0)
		default:
			executeCommand(command, args, currentDir)
		}
	}
}

func startWebServer() {
    http.HandleFunc("/ws", handleConnections)
    log.Println("http server started on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func parseCommand(input string) (string, string) {
	parts := strings.SplitN(input, " ", 2)
	var command, args string
	command = parts[0]
	if len(parts) > 1 {
		args = parts[1]
	}
	return command, args
}

func changeDirectory(path string, currentDir string) string {
	if path == "" {
		fmt.Println("cd: path required")
		return currentDir
	}
	err := os.Chdir(path)
	if err != nil {
		fmt.Println("cd:", err)
		return currentDir
	}
	newDir, err := os.Getwd()
	if err != nil {
		fmt.Println("cd:", err)
		return currentDir
	}
	return newDir
}

func listFiles(path string) {
	if path == "" {
		path = "."
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("ls:", err)
		return
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func catFile(filename string) {
	if filename == "" {
		fmt.Println("cat: filename required")
		return
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("cat:", err)
		return
	}
	fmt.Println(string(content))
}

func showHistory(history []string) {
	for i, command := range history {
		fmt.Printf("%d %s\n", i, command)
	}
}

func makeDirectory(path string) {
	if path == "" {
		fmt.Println("mkdir: path required")
		return
	}
	err := os.Mkdir(path, 0750)
	if err != nil {
		fmt.Println("mkdir:", err)
	}
}

func removeFile(path string) {
	if path == "" {
		fmt.Println("rm: path required")
		return
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println("rm:", err)
	}
}

func executeCommand(command string, args string, dir string) {
	cmd := exec.Command(command, args)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))
}

