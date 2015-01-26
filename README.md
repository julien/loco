lr
==

An example http/livereload server
with [Go](http://golang.org)

Usage
=====

```
Usage of lr:
  -cache=30: number of days for cache/expires header
  -excludes="": directories to exclude when watching
  -port="8000": default port
  -recursive=false: watch for file changes in all directories
  -root=".": root directory
```

+ Install with:

  `go install github.com/julien/lr`


+ Navigate to a directory you want to use as the "root" directory:

  `cd ~/somedir # for example`

+ Start the server:

  `lr -port 8000 -root . -recursive -excludes=node_modules,bower_components # flags are optional`

+ If you want to be notified about file changes
  include this `script` tag:

  ```html
  <!-- NOTE: change 8000 to the port you used to start the server -->
  <script src="http://localhost:8000/livereload.js"></script>
  ```

+ Check the [report card](http://goreportcard.com/report/julien/lr)
