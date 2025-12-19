// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// Package errors provides simple functions to manipulate errors.
//
// Very note that, this package is quite a base package, which should not import extra
// packages except standard packages, to avoid cycle imports.
package gerror

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// apiCode is the interface for Code feature.
type apiCode interface {
	Error() string // It should be an error.
	Code() int
}

// apiStack is the interface for Stack feature.
type apiStack interface {
	Error() string // It should be an error.
	StackString() string
}

// apiCause is the interface for Cause feature.
type apiCause interface {
	Error() string // It should be an error.
	Cause() error
}

// apiCurrent is the interface for Current feature.
type apiCurrent interface {
	Error() string // It should be an error.
	Current() error
}

// apiNext is the interface for Next feature.
type apiNext interface {
	Error() string // It should be an error.
	Next() error
}

// New creates and returns an error which is formatted from given msg.
func New(text string) error {
	return &Error{
		stack: callers(),
		msg:   text,
		code:  -1,
	}
}

// Newf returns an error that formats as the given format and args.
func Newf(format string, args ...interface{}) error {
	return &Error{
		stack: callers(),
		msg:   fmt.Sprintf(format, args...),
		code:  -1,
	}
}

// NewSkip creates and returns an error which is formatted from given msg.
// The parameter <skip> specifies the stack callers skipped amount.
func NewSkip(skip int, text string) error {
	return &Error{
		stack: callers(skip),
		msg:   text,
		code:  -1,
	}
}

// NewSkipf returns an error that formats as the given format and args.
// The parameter <skip> specifies the stack callers skipped amount.
func NewSkipf(skip int, format string, args ...interface{}) error {
	return &Error{
		stack: callers(skip),
		msg:   fmt.Sprintf(format, args...),
		code:  -1,
	}
}

// Wrap wraps error with msg.
// It returns nil if given err is nil.
func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		msg:   text,
		code:  Code(err),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// It returns nil if given <err> is nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		msg:   fmt.Sprintf(format, args...),
		code:  Code(err),
	}
}

// WrapSkip wraps error with msg.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapSkip(skip int, err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		msg:   text,
		code:  Code(err),
	}
}

// WrapSkipf wraps error with msg that is formatted with given format and args.
// It returns nil if given err is nil.
// The parameter <skip> specifies the stack callers skipped amount.
func WrapSkipf(skip int, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(skip),
		msg:   fmt.Sprintf(format, args...),
		code:  Code(err),
	}
}

// NewCode creates and returns an error that has error code and given msg.
func NewCode(code int, text string) Error {
	return Error{
		stack: callers(),
		msg:   text,
		code:  code,
	}
}
func NewCodeTest(code int, text string, s []uintptr) Error {
	return Error{
		stack: s,
		msg:   text,
		code:  code,
	}
}

// Cause returns the error code of current error.
// It returns -1 if it has no error code or it does not implements interface Code.
func Code(err error) int {
	if err != nil {
		if e, ok := err.(apiCode); ok {
			return e.Code()
		}
	}
	return -1
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return true
	}
	msg := err.Error()
	if strings.Contains(msg, "i/o timeout") || strings.Contains(msg, "connection reset by peer") {
		return true
	}
	return false
}
