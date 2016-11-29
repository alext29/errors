package errors

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// emsg is a helper struct that holds information for a single error
// in an error stack.
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

// Error represents a single error stack.
type Error struct {
	errs []*emsg
}

// Depth of error stack.
func (err *Error) Depth() int {
	return len(err.errs)
}

// Header information for an error in error stack.
func (err *Error) Header(i int) (string, int, time.Time) {
	if i < 0 || i >= err.Depth() {
		return "", -1, time.Time{}
	}
	return err.errs[i].header()
}

// Msg for an error in error stack.
func (err *Error) Msg(i int) string {
	if i < 0 || i >= err.Depth() {
		return ""
	}
	return err.errs[i].msg()
}

// New creates a new error.
func New(format string, args ...interface{}) error {
	file, line := fileline(2)
	err := &Error{}
	return err.add(file, line, fmt.Errorf(format, args...))
}

// Wrap adds a error to stack with additional context.
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
		return err.errs[0].e
	}
	return e
}

// add is a helper function that adds an error to error stack.
func (err *Error) add(file string, line int, e error) *Error {
	m := &emsg{
		time: time.Now(),
		file: file,
		line: line,
		e:    e,
	}
	err.errs = append(err.errs, m)
	return err
}

// Error implements error interface.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	first := true
	var str string
	for i := 0; i < len(err.errs); i++ {
		file, line, _ := err.Header(i)
		if !first {
			str = fmt.Sprintf("%s :: ", str)
		}
		first = false
		str = fmt.Sprintf("%s%s:%d %s", str, file, line, err.Msg(i))
	}
	return str
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
