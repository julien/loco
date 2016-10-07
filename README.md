loco
====

![https://circleci.com/gh/julien/loco](https://circleci.com/gh/julien/loco.png?circle-token=722cb47155b6d2b3983203591655815031c46b08)


An basic http server with [Go](http://golang.org)

Usage
=====

```shell
loco -port=PORT -root=ROOT_DIRECTORY
```

```shell
Usage of loco:
  -port="8000": default port
  -root=".": root directory
```

+ Install with:

  `go install github.com/julien/loco`


+ Navigate to a directory you want to use as the "root" directory:

  `cd ~/somedir # for example`

+ Start the server:

  `loco -port 8000 -root .`

+ Check the [report card](http://goreportcard.com/report/julien/loco)



