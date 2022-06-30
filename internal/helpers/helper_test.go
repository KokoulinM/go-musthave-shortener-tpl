package helpers

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Simple test #1",
			args: args{len: 5},
			want: 6,
		},
		{
			name: "Negative test",
			args: args{len: 5},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomString(tt.args.len); len(got) == tt.want {
				t.Errorf("RandomString() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
