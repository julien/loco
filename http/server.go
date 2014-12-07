package http

import (
	"github.com/gorilla/websocket"
	"gopkg.in/fsnotify.v1"
	"log"
	"net/http"
	"path"
	"regexp"
)

const DEFAULT_PORT string = "8000"

var (
	valid    *regexp.Regexp     = regexp.MustCompile(`\d{4}`)
	upgrader websocket.Upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	watcher  *fsnotify.Watcher
)

func Serve(port, root string) {

	if !valid.MatchString(port) {
		port = DEFAULT_PORT
	}

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Watcher create error %s\n", err)
	}
	defer watcher.Close()

	err = watcher.Add(root)
	if err != nil {
		log.Printf("Watcher add error %s\n", err)
	}

	log.Println("Watching for file changes")
	log.Printf("Starting server: 0.0.0.0:%s - Root directory: %s\n", port, path.Dir(root))

	http.HandleFunc("/ws", socket)
	http.Handle("/", http.FileServer(http.Dir(root)))
	http.ListenAndServe(":"+port, nil)
}

func socket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error %s\n", err)
		return
	}
	go writer(c)
	// reader(c)
}

// func reader(c *websocket.Conn) {
// 	for {
// 		_, _, err := c.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		// fmt.Println("Message Type", msgType, "Body: ", string(p[0:]))
// 	}
// }

func writer(c *websocket.Conn) {
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("Modified file: %s\n", ev.Name)
				if err := c.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
					return
				}
			}
		case err := <-watcher.Errors:
			log.Printf("Watch error %s:\n", err)
		}
	}
}
