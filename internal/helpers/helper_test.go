package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGenerateRandom(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{
				size: 10,
			},
			want: []byte{0x52, 0xfd, 0xfc, 0x7, 0x21, 0x82, 0x65, 0x4f, 0x16, 0x3f},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GenerateRandom(tt.args.size)

			assert.Equal(t, tt.want, got)
		})
	}
}
