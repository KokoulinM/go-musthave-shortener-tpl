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
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShorterURL(tt.args.longURL); got != tt.want {
				t.Errorf("ShorterURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
