package entity

import (
	"gohive/internal/crypt"

	log "github.com/Sirupsen/logrus"
)

var (
	PRIVATE_KEY string
	PUBLIC_KEY  string
)

func init() {
	privateKey, pubKey, err := crypt.GenRSAKeyPair(2048)
	if err != nil {
		log.Fatal("gen key pair failed: ", err)
	} else {
		PRIVATE_KEY = privateKey
		PUBLIC_KEY = pubKey
	}
}
