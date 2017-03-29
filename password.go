package dbase

import (
	"crypto/rand"

	"golang.org/x/crypto/scrypt"
)

type Password struct {
	Hash, Salt []byte
}

func NewPassword(pw string) (Password, error) {
	salt := make([]byte, 10)
	_, err := rand.Read(salt)
	if err != nil {
		return Password{}, err
	}
	dk, err := scrypt.Key([]byte(pw), salt, 16384, 8, 1, 32)
	return Password{dk, salt}, nil

}
func (p Password) Check(pw string) bool {
	if pw == "" {
		return false
	}
	dk, err := scrypt.Key([]byte(pw), p.Salt, 16384, 8, 1, 32)
	if err != nil {
		return false
	}
	return string(dk) == string(p.Hash)
}
