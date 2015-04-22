loco
====

![Buidl status](https://circleci.com/gh/julien/loco.png?circle-token=722cb47155b6d2b3983203591655815031c46b08)


An example http/livereload server
with [Go](http://golang.org)

Usage
=====

```shell
loco -port=PORT -root=ROOT_DIRECTORY FILES
```

```shell
Usage of loco:
  -port="3000": default port
  -root=".": root directory
```

+ Install with:

  `go install github.com/julien/loco`


+ Navigate to a directory you want to use as the "root" directory:

  `cd ~/somedir # for example`

+ Start the server:

  `loco -port 8000 -root . *.js # files are optional`

+ If you want to be notified about file changes
  include this `script` tag:

  ```html
  <!-- NOTE: change 8000 to the port you used to start the server -->
  <script src="http://localhost:8000/livereload.js"></script>
  ```

+ Check the [report card](http://goreportcard.com/report/julien/loco)
