package main

import (
	"flag"
	"github.com/julien/lr/http"
)

var (
	port string
	root string
)

func main() {
	flag.StringVar(&port, "port", "8000", "default port")
	flag.StringVar(&root, "root", ".", "root directory")
	flag.Parse()

	http.Start(port, root)
}
