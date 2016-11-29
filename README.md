# Simplified golang errors
[![GoDoc](https://godoc.org/github.com/alext29/errors?status.svg)](https://godoc.org/github.com/alext29/errors)

This package provides methods to handle and help log golang errors. Use this with the modified glog package to super simplify error handling in your code.

## Install

```bash
go get github.com/alext29/glog
go get github.com/alext29/errors
```
The glog package patches the original golang/glog to handle errors.Error type.

Logging error captures the entire error stack correctly, including line numbers and timestamps. And prints the error stack on separate lines, providing information similar to a stack trace.

## Examples

The code exposes the following functions
```go
New(format string, args ...interface{}) error
```
Creates a new error.

```go
Wrap(e error, format string, args ...interface{}) error
```
Wrap adds the additional information to the error stack.

```go
Cause(e error) error
```
Cause returns the first error.

## Examples

### Example 1

```go
top.go
1  func topFn(){
2      arg := "test"
3      err := midFn(arg)
4      if err != nil {
5          err = errors.Wrap(err, "failed mid function: %s", arg)
6          glog.Error(err)
7      }
8  }

mid.go
1  func midFn(arg string) error {
2      err := extFn()
3      if err != nil {
4          return errors.Wrap(err, "external function failed")
5       }
6       return nil
7  }
8
9  func extFn() error {
10     return fmt.Errorf("external error")
11 }
```

This code will generate log smilar to the one below. Note the filename, line numbers and timestamps are accurately recorded.

```bash
E1128 21:19:07.663216    3586 mid.go:4] external error
E1128 21:19:07.663222    3586 mid.go:4] external function failed
E1128 21:19:07.663226    3586 top.go:5] failed mid function: test
```

### Example 2

```go
top.go
1  func topFn(){
2      arg := "test"
3      err := midFn(arg)
4      if err != nil {
5          err = errors.Wrap(err, "failed mid function: %s", arg)
6          glog.Error(err)
7      }
8  }

mid.go
1  func midFn(arg string) error {
3      if errorCondition {
4          return errors.New("reached an error condition: %s", arg)
5      }
6      return nil
7  }
```

This code will generate log smilar to the one below. Note the filename, line numbers and timestamps are accurately recorded.

```bash
E1128 21:19:07.663222    3586 mid.go:4] reached error condition: test
E1128 21:19:07.663226    3586 top.go:5] failed mid function: test
```
