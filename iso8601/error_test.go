package iso8601

import "testing"

func TestUnexpectedTokenError_Error(t *testing.T) {
	tests := []struct {
		err  *UnexpectedTokenError
		want string
	}{
		{
			err: &UnexpectedTokenError{
				Value:      "2020/10/01",
				Token:      "/",
				AfterToken: "2020",
				Expected:   "-",
			},
			want: `unexpected token "/" after "2020" expected "-" ("2020/10/01")`,
		},
		{
			err: &UnexpectedTokenError{
				Value:    "2020/10/01",
				Token:    "/",
				Expected: "-",
			},
			want: `unexpected token "/" expected "-" ("2020/10/01")`,
		},
		{
			err: &UnexpectedTokenError{
				Value: "2020/10/01",
				Token: "/",
			},
			want: `unexpected token "/" ("2020/10/01")`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("UnexpectedTokenError.Error()\n+ %q\n- %q", got, tt.want)
			}
		})
	}
}
