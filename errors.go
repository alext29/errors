package errors

// Package errors integrates with glog and provides a super
// simplified way for handling and logging errors.

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// New creates a new error.
func New(format string, args ...interface{}) error {
	file, line := fileline(2)
	err := &Error{}
	return err.add(file, line, fmt.Errorf(format, args...))
}

// Wrap adds error to stack with additional context.
func Wrap(e error, format string, args ...interface{}) error {
	if e == nil {
		return nil
	}
	file, line := fileline(2)
	err, ok := e.(*Error)
	if !ok {
		err = &Error{}
		err = err.add(file, line, e)
	}
	return err.add(file, line, fmt.Errorf(format, args...))
}

// Cause returns the root error.
func Cause(e error) error {
	if e == nil {
		return nil
	}
	if err, ok := e.(*Error); ok {
		return err.stack[0].e
	}
	return e
}

// Error represents an error stack.
type Error struct {
	stack []*emsg
}

// Depth of error stack.
func (err *Error) Depth() int {
	return len(err.stack)
}

// Header information for ith error in stack.
func (err *Error) Header(i int) (string, int, time.Time) {
	if i < 0 || i >= err.Depth() {
		return "", -1, time.Time{}
	}
	return err.stack[i].header()
}

// Msg for ith error in stack.
func (err *Error) Msg(i int) string {
	if i < 0 || i >= err.Depth() {
		return ""
	}
	return err.stack[i].msg()
}

// Error implements error interface.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	first := true
	var str string
	for i := 0; i < err.Depth(); i++ {
		if !first {
			str = fmt.Sprintf("%s :: ", str)
		}
		first = false
		file, line, _ := err.Header(i)
		str = fmt.Sprintf("%s%s:%d %s", str, file, line, err.Msg(i))
	}
	return str
}

// add an error to error stack.
func (err *Error) add(file string, line int, e error) *Error {
	err.stack = append(err.stack, &emsg{time.Now(), file, line, e})
	return err
}

// emsg is single entry in the error stack.
type emsg struct {
	time time.Time
	file string
	line int
	e    error
}

// header for an emsg.
func (m *emsg) header() (string, int, time.Time) {
	return m.file, m.line, m.time
}

// msg for an emsg.
func (m *emsg) msg() string {
	return m.e.Error()
}

// fileline retrieves the filename and line for calling function at given depth.
func fileline(depth int) (string, int) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return file, line
}
