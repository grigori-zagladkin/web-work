package slug

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"latin", "Hello, World!", "hello-world"},
		{"cyrillic", "Привет, Go!", "privet-go"},
		{"hard sign", "Объект", "obekt"},
		{"punctuation and spaces", "  --Go 1.25 & Generics--  ", "go-1-25-generics"},
		{"empty", "", ""},
		{"only symbols", "!!!", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Make(tt.in); got != tt.want {
				t.Errorf("Make(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
