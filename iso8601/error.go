package iso8601

import (
	"errors"
	"fmt"
	"strings"
)

// UnexpectedTokenError represents an error encountered when an unexpected
// token is detected during parsing. This error provides details about
// the token that was unexpected, any preceding token, and what was
// expected in its place.
type UnexpectedTokenError struct {
	Value      string
	Token      string
	AfterToken string
	Expected   string
}

// Error implements the error interface.
func (u *UnexpectedTokenError) Error() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "unexpected token %q", u.Token)
	if u.AfterToken != "" {
		fmt.Fprintf(&buf, " after %q", u.AfterToken)
	}
	if u.Expected != "" {
		fmt.Fprintf(&buf, " expected %q", u.Expected)
	}
	fmt.Fprintf(&buf, " (%q)", u.Value)
	return buf.String()
}

func overrideUnexpectedTokenValue(err error, b []byte) error {
	var unexpected *UnexpectedTokenError
	if errors.As(err, &unexpected) {
		unexpected.Value = string(b)
		err = unexpected
	}
	return err
}
