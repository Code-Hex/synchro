package iso8601

import (
	"fmt"
	"strings"
)

type UnexpectedTokenError struct {
	Value      string
	Token      string
	AfterToken string
	Expected   string
}

var _ error = (*UnexpectedTokenError)(nil)

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
