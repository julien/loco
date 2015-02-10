package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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
	// cache     int
	recursive bool
	excludes  string
)

func init() {
	flag.StringVar(&port, "port", "3000", "default port")
	flag.StringVar(&root, "root", ".", "root directory")
	// flag.IntVar(&cache, "cache", 30, "number of days for cache/expires header")
	flag.BoolVar(&recursive, "recursive", false, "watch for file changes in all directories")
	flag.StringVar(&excludes, "excludes", "", "directories to exclude when watching")
}

func main() {
	flag.Parse()
	checkPort(port)

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Watcher create error %s\n", err)
	}
	fmt.Println("Watching for file changes")

	fileChan := make(chan string)
	go addFiles(root, fileChan)

	defer watcher.Close()

	fmt.Printf("Starting server: 0.0.0.0:%s - Root directory: %s\n", port, path.Dir(root))

	http.Handle("/", gzHandler(fileHandler(root)))
	http.Handle("/ws", socketHandler())
	http.Handle("/livereload.js", scriptHandler())

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func checkPort(port string) string {
	if !valid.MatchString(port) {
		port = defaultport
	}
	return port
}

func addFiles(root string, fileChan chan string) {
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
		fileChan <- files[i]
		fmt.Printf("Added %s to watcher\n", files[i])
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
			fmt.Printf("Upgdrader error %s\n", err)
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
				fmt.Printf("Modified file: %s\n", ev.Name)
				for cl := range clients {
					if err := cl.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
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
			if _, b := clients[c]; b {
				clients[c] = false
				delete(clients, c)
			}
			break
		}
	}
}

// func cacheHandler(days int, next http.Handler) http.Handler {
//
// 	if days < 1 {
// 		days = 1
// 	}
// 	age := days * 24 * 60 * 60 * 1000
// 	t := time.Now().Add(time.Duration(time.Hour * 24 * time.Duration(days)))
//
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
// 		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(age))
// 		w.Header().Set("Expires", t.Format(time.RFC1123Z))
//
// 		next.ServeHTTP(w, r)
// 	})
// }

type gzResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gw := gzResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gw, r)

	})
}
