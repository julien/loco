package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/fsnotify.v1"
)

const (
	defaultport = "3000"
	maxfiles    = 300
	script      = `(function () { window.addEventListener('load', function () {
  var ws = new WebSocket('ws://localhost:%s/ws');
  ws.onmessage = function () { ws.close(); location = location; };
});}());`
)

var (
	valid    = regexp.MustCompile(`\d{4}`)
	hidden   = regexp.MustCompile(`^\.`)
	upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	clients  = make(map[*websocket.Conn]bool)
	watcher  *fsnotify.Watcher
	port     string
	root     string
)

func init() {
	flag.StringVar(&port, "port", "3000", "default port")
	flag.StringVar(&root, "root", ".", "root directory")
}

func main() {
	flag.Parse()
	checkPort(port)

	globs := flag.Args()
	if len(globs) > 0 {
		var files []string
		for _, glob := range globs {

			matches, err := filepath.Glob(glob)
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, matches...)
		}

		if len(files) > 0 {
			var err error
			watcher, err = fsnotify.NewWatcher()
			if err != nil {
				fmt.Printf("could not create watcher\n", err)
			}
			fmt.Println("watching for file changes")

			fileChan := make(chan string)
			go addFiles(files, fileChan)
			defer watcher.Close()

			http.Handle("/ws", socketHandler())
			http.Handle("/livereload.js", scriptHandler())
		}
	}

	fmt.Printf("starting server: 0.0.0.0:%s - root directory: %s\n", port, path.Dir(root))
	http.Handle("/", gzHandler(fileHandler(root)))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func checkPort(port string) string {
	if !valid.MatchString(port) {
		port = defaultport
	}
	return port
}

func addFiles(files []string, fileChan chan string) {
	max := len(files)
	if len(files) > maxfiles {
		max = maxfiles
	}

	for i := 0; i < max; i++ {
		time.Sleep(10 * time.Millisecond)

		if err := watcher.Add(files[i]); err != nil {
			return
		}
		fileChan <- files[i]
		fmt.Printf("added %s to watcher\n", files[i])
	}
}

func fileHandler(root string) http.Handler {
	return http.FileServer(http.Dir(root))
}

func scriptHandler() http.Handler {
	s := []byte(fmt.Sprintf(script, port))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(s)
	})
}

func socketHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("websocket error %s\n", err)
			return
		}

		clients[c] = true
		go writer(c)
		reader(c)
	})
}

func writer(c *websocket.Conn) {
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf("modified file: %s\n", ev.Name)
				for cl := range clients {
					if err := cl.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
						return
					}
				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("watch error %s:\n", err)
		}
	}
}

func reader(c *websocket.Conn) {
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if _, b := clients[c]; b {
				clients[c] = false
				delete(clients, c)
			}
			break
		}
	}
}
