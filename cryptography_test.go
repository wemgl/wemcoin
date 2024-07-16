package wemcoin

import (
	"testing"
)

func TestSHA256hash(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generates hash for string",
			args: args{s: "Hello world!"},
			want: "c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a",
		},
		{
			name: "generates hash for empty string",
			args: args{s: ""},
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA256hash(tt.args.s); got != tt.want {
				t.Errorf("SHA256hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenKey(t *testing.T) {
	got := GenKey()
	if got == nil {
		t.Errorf("GenKey() result was nil")
	}

	if got.Public() == nil {
		t.Errorf("public key was nil")
	}
}

func TestSign(t *testing.T) {
	key := GenKey()
	h := SHA256hash("Hello world!")
	if got := Sign(key, h); len(got) < 0 {
		t.Errorf("Sign(%v, %s) = len(%v), want = >0", key, h, len(got))
	}
}

func TestVerify(t *testing.T) {
	key := GenKey()
	h := SHA256hash("Hello world!")
	signed := Sign(key, h)
	if got := Verify(&key.PublicKey, h, signed); !got {
		t.Errorf("Verify(%v, %s, %v) = %t, want = %t", key.PublicKey, h, signed, got, true)
	}
}
