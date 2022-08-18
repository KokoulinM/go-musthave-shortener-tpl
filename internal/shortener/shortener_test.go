package shortener

import "testing"

func TestShorterURL(t *testing.T) {
	type args struct {
		longURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test #1",
			args: args{
				longURL: "123",
			},
			want: "QL0AFWMIX8NRZTKeof9cXsvbvu8=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShorterURL(tt.args.longURL); got != tt.want {
				t.Errorf("ShorterURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
