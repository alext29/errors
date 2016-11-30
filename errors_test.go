package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	msgs := []string{
		"error cause",
		"error wrapper #1",
		"error wrapper #2",
		"error wrapper #3",
	}

	depth := len(msgs)
	err := New(msgs[0])
	for ii := 1; ii < depth; ii++ {
		Wrap(err, msgs[ii])
	}
	for ii, msg := range msgs {
		e := err.(*Error)
		assert.Equal(t, msg, e.Msg(ii))
	}
	assert.Equal(t, msgs[0], Cause(err).Error())
}

func TestNil(t *testing.T) {
	var err error
	e := Wrap(err, "adding to a nil error")
	assert.Nil(t, e)
	assert.Nil(t, Cause(e))
}

func TestNew(t *testing.T) {
	tests := []string{
		"error string",
		"",
	}

	for _, test := range tests {
		err := New(test)
		assert.Equal(t, test, Cause(err).Error())
	}
}
