# Simplified golang error handling and logging.
[![GoDoc](https://godoc.org/github.com/alext29/errors?status.svg)](https://godoc.org/github.com/alext29/errors)

This errors package integrates with glog and provides a super simplified way for handling and logging errors.

## Install

```bash
go get github.com/alext29/glog
go get github.com/alext29/errors
```
The glog package patches the original golang/glog to work with errors package.

## Error Handling

This package exposes the following methods to handle errors.

Create a new error.
```go
New(format string, args ...interface{}) error
```

Wrap adds the additional information to the error stack.
```go
Wrap(e error, format string, args ...interface{}) error
```

Cause returns the first error in error stack.
```go
Cause(e error) error
```
Additional exposed methods are used by glog package to log error trace.

## Error Logging

Log errors using glog primitives Info(error), Warning(error), Error(error) or Fatal(error). This would log the relevant error stack along with line numbers and file names.

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
Expected log output.
```bash
E1128 21:19:07.663216    3586 mid.go:4] external error
E1128 21:19:07.663222    3586 mid.go:4] external function failed
E1128 21:19:07.663226    3586 top.go:5] failed mid function: test
```

### Example 2

Same example with some idiomatic simplification.

```go
top.go
1 func topFn(){
2     arg := "test"
3     err := midFn(arg)
4     if err != nil {
5         err = errors.Wrap(err, "failed mid function: %s", arg)
6         glog.Error(err)
7     }
8 }

mid.go
1 func midFn(arg string) error {
2     err := extFn()
3     return errors.Wrap(err, "external function failed")
4 }
5
6 func extFn() error {
7     return fmt.Errorf("external error")
8 }
```
Expected log output.
```bash
E1128 21:19:07.663216    3586 mid.go:3] external error
E1128 21:19:07.663222    3586 mid.go:3] external function failed
E1128 21:19:07.663226    3586 top.go:5] failed mid function: test
```

### Example 3

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
Expected log output.
```bash
E1128 21:19:07.663222    3586 mid.go:4] reached error condition: test
E1128 21:19:07.663226    3586 top.go:5] failed mid function: test
```
