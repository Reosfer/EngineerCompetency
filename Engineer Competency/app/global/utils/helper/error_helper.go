package helper

import (
	"github.com/sirupsen/logrus"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func Recover(location string) {
	if r := recover(); r != nil {
		logrus.Debugf("recover panic action from %s : %s\n", location, r)
	}
}

// New returns an error that formats as the given text.
func NewError(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
