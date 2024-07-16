package wemcoin

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func SHA256hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GenKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return key
}

func Sign(p *ecdsa.PrivateKey, h string) []byte {
	sig, err := p.Sign(rand.Reader, []byte(h), crypto.SHA256)
	if err != nil {
		panic(err)
	}
	return sig
}

func Verify(p *ecdsa.PublicKey, d string, s []byte) bool {
	return ecdsa.VerifyASN1(p, []byte(d), s)
}
