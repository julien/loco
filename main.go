package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fsnotify.v1"
	"net/http"
	"path"
	"regexp"
)

const defaultport string = "8000"

var (
	valid    = regexp.MustCompile(`\d{4}`)
	upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	watcher  *fsnotify.Watcher
	port     string
	root     string
)

func main() {
	flag.StringVar(&port, "port", "8000", "default port")
	flag.StringVar(&root, "root", ".", "root directory")
	flag.Parse()

	if !valid.MatchString(port) {
		port = defaultport
	}

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Watcher create error %s\n", err)
	}
	defer watcher.Close()

	err = watcher.Add(root)
	if err != nil {
		fmt.Printf("Watcher add error %s\n", err)
	}

	fmt.Println("Watching for file changes")
	fmt.Printf("Starting server: 0.0.0.0:%s - Root directory: %s\n", port, path.Dir(root))

	http.HandleFunc("/livereload.js", livereload)
	http.HandleFunc("/ws", socket)
	http.Handle("/", http.FileServer(http.Dir(root)))
	http.ListenAndServe(":"+port, nil)
}

func livereload(w http.ResponseWriter, r *http.Request) {
	script := `(function () {
  var ws = new WebSocket('ws://localhost:%s/ws');
  ws.onmessage = function () { document.location.reload(); };
}())
    `
	s := []byte(fmt.Sprintf(script, port))
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(s)
}

func socket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	go writer(c)
}

func writer(c *websocket.Conn) {
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf("Modified file: %s\n", ev.Name)
				if err := c.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
					return
				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("Watch error %s:\n", err)
		}
	}
}
