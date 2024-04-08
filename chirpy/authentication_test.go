package main

import (
	"testing"
	"time"
)

func Test_createSignedString(t *testing.T) {
	type args struct {
		id       int
		issuer   string
		duration time.Duration
		secret   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Generate a signed string",
			args: args{
				id:       1,
				issuer:   "chirpy-access",
				duration: 5 * time.Second,
				secret:   "sausages",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createSignedString(
				tt.args.id,
				tt.args.issuer,
				tt.args.duration,
				tt.args.secret,
			)
			t.Logf("signed string is:\n%s", got)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"createSignedString() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
		})
	}
}
