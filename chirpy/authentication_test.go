package main

import (
	"testing"
)

func Test_createSignedString(t *testing.T) {
	type args struct {
		id            int
		email         string
		expiresInSecs int
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
				id:            1,
				expiresInSecs: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createSignedString(
				tt.args.id,
				tt.args.expiresInSecs,
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
