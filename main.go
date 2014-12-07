package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fsnotify.v1"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	defaultport string = "8000"
	maxfiles           = 30
)

var (
	valid     = regexp.MustCompile(`\d{4}`)
	hidden    = regexp.MustCompile(`^\.`)
	upgrader  = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	clients   = make(map[*websocket.Conn]bool)
	watcher   *fsnotify.Watcher
	port      string
	root      string
	recursive bool
	excludes  string
)

func main() {
	flag.StringVar(&port, "port", "8000", "default port")
	flag.StringVar(&root, "root", ".", "root directory")
	flag.BoolVar(&recursive, "recursive", false, "watch for file changes in all directories")
	flag.StringVar(&excludes, "excludes", "", "directories to exclude when watching")
	flag.Parse()

	if !valid.MatchString(port) {
		port = defaultport
	}

	fmt.Println(excludes)

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Watcher create error %s\n", err)
	}
	fmt.Println("Watching for file changes")
	defer watcher.Close()
	add(root)

	fmt.Printf("Starting server: 0.0.0.0:%s - Root directory: %s\n", port, path.Dir(root))

	http.HandleFunc("/livereload.js", livereload)
	http.HandleFunc("/ws", socket)
	http.Handle("/", http.FileServer(http.Dir(root)))
	http.ListenAndServe(":"+port, nil)
}

func add(root string) {
	var files []string
	files = append(files, root)

	if recursive == true {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && !hidden.MatchString(path) && !strings.Contains(excludes, path) {
				files = append(files, path)
			}
			return nil
		})
	}

	max := len(files)
	if len(files) > maxfiles {
		max = maxfiles
	}

	for i := 0; i < max; i++ {
		time.Sleep(10 * time.Millisecond)

		if err := watcher.Add(files[i]); err != nil {
			fmt.Printf("Watcher add error %s\n", err)
			return
		}
		fmt.Printf("Added %s to watcher\n", files[i])
	}
}

func livereload(w http.ResponseWriter, r *http.Request) {
	script := `(function () {
  window.onload = function () {
    var ws = new WebSocket('ws://localhost:%s/ws');
    ws.onmessage = function () {
      ws.close();
      location = location;
    };
  };
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

	clients[c] = true
	go writer(c)
	reader(c)
}

func writer(c *websocket.Conn) {
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf("Modified file: %s\n", ev.Name)

				for cl := range clients {
					// fmt.Println("C", b)
					if err := cl.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
						fmt.Printf("Error writing message: %s\n", err)
						return
					}

				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("Watch error %s:\n", err)
		}
	}
}

func reader(c *websocket.Conn) {
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("ReadMessage error: %s\n", err)

			if _, b := clients[c]; b {
				clients[c] = false
				delete(clients, c)
			}

			break
		}
	}
}
