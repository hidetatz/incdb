package main

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	eof := &Token{Type: TkEOF}
	testcases := []struct {
		query string
		want  *Token
	}{
		{
			query: "r",
			want:  &Token{Type: TkRead, Next: eof},
		},
		{
			query: "   r   ",
			want:  &Token{Type: TkRead, Next: eof},
		},
		{
			query: "   r   123",
			want: &Token{
				Type: TkRead,
				Next: &Token{
					Type: TkStr,
					Val:  "123",
					Next: eof,
				},
			},
		},
		{
			query: "   w   abc   456   ",
			want: &Token{
				Type: TkWrite,
				Next: &Token{
					Type: TkStr,
					Val:  "abc",
					Next: &Token{
						Type: TkStr,
						Val:  "456",
						Next: eof,
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.query, func(t *testing.T) {
			got := tokenize(tc.query)
			for true {
				if got == nil {
					t.Fatalf("got is unexpectedly nil while expected is: %v", tc.want)
				}

				if tc.want == nil {
					t.Fatalf("want is unexpectedly nil while got is: %v", got)
				}

				if got.Type != tc.want.Type || got.Val != tc.want.Val {
					t.Fatalf("got: %v, expected: %v", got, tc.want)
				}

				if got.Next == nil && tc.want.Next == nil {
					break
				}

				got = got.Next
				tc.want = tc.want.Next
			}
		})
	}
}
