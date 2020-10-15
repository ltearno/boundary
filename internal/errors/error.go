package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// Op represents an operation (package.function).
// For example iam.CreateRole
type Op string

// Err provides the ability to specify a Msg, Op, Code and Wrapped error.
// Errs must have a Code and all other fields are optional. We've chosen Err
// over Error for the identifier to support the easy embedding of Errs.  Errs
// can be embedded without a conflict between the embedded Err and Err.Error().
type Err struct {
	// Code is the error's code, which can be used to get the error's
	// errorCodeInfo, which contains the error's Kind and Message
	Code Code

	// Msg for the error
	Msg string

	// Op represents the operation raising/propagating an error and is optional
	Op Op

	// Wrapped is the error which this Error wraps and will be nil if there's no
	// error to wrap.
	Wrapped error
}

// New creates a new Error and supports the options of:
// WithMsg() - allows you to specify an optional error msg, if the default
// msg for the error Code is not sufficient.
// WithWrap() - allows you to specify
// an error to wrap
func New(c Code, opt ...Option) error {
	opts := GetOpts(opt...)
	return &Err{
		Code:    c,
		Wrapped: opts.withErrWrapped,
		Msg:     opts.withErrMsg,
	}
}

/// Convert will convert the error to a Boundary Error and attempt to add a
//helpful error msg as well. If that's not possible, it return nil
func Convert(e error) error {
	// nothing to convert.
	if e == nil {
		return nil
	}

	var alreadyConverted *Err
	if errors.As(e, &alreadyConverted) {
		return alreadyConverted
	}

	var pqError *pq.Error
	if errors.As(e, &pqError) {
		if pqError.Code.Name() == "unique_violation" {
			return New(NotUnique, WithMsg(pqError.Detail), WithWrap(ErrNotUnique))
		}
		if pqError.Code.Name() == "not_null_violation" {
			msg := fmt.Sprintf("%s must not be empty", pqError.Column)
			return New(NotNull, WithMsg(msg), WithWrap(ErrNotNull))
		}
		if pqError.Code.Name() == "check_violation" {
			msg := fmt.Sprintf("%s constraint failed", pqError.Constraint)
			return New(CheckConstraint, WithMsg(msg), WithWrap(ErrCheckConstraint))
		}
	}
	// unfortunately, we can't help.
	return e
}

// Info about the Error
func (e *Err) Info() Info {
	if info, ok := errorCodeInfo[e.Code]; ok {
		return info
	}
	return errorCodeInfo[Unknown]
}

// Error satisfies the error interface and returns a string representation of
// the error.
func (e *Err) Error() string {
	var s strings.Builder
	if e.Op != "" {
		join(&s, ": ", string(e.Op))
	}
	if e.Msg != "" {
		join(&s, ": ", e.Msg)
	}

	if info, ok := errorCodeInfo[e.Code]; ok {
		if e.Msg == "" {
			join(&s, ": ", info.Message) // provide a default.
			join(&s, ", ", info.Kind.String())
		} else {
			join(&s, ": ", info.Kind.String())
		}
	}
	join(&s, ": ", fmt.Sprintf("error #%d", e.Code))

	if e.Wrapped != nil {
		join(&s, ": \n", e.Wrapped.Error())
	}
	return s.String()
}

func join(str *strings.Builder, delim string, s string) {
	if str.Len() == 0 {
		_, _ = str.WriteString(s)
		return
	}
	_, _ = str.WriteString(delim + s)
}

// Unwrap implements the errors.Unwrap interface and allows callers to use the
// errors.Is() and errors.As() functions effectively for any wrapped errors.
func (e *Err) Unwrap() error {
	return e.Wrapped
}