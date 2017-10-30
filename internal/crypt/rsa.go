package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	//"encoding/base64"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	//"io/ioutil"
)

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) ([]byte, error) {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	return asn1.Marshal(info)
}

func RSASign(privateKey string, hash crypto.Hash, msg []byte) ([]byte, error) {
	p, _ := pem.Decode([]byte(privateKey))
	if p == nil {
		return nil, errors.New("decode pem failed...")
	}

	//key, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	key, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, err
	}

	h := hash.New()
	h.Write(msg)
	hashed := h.Sum(nil)
	//sig, err := rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), hash, hashed[:])
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, hash, hashed[:])
	return sig, err
}

func RSAVerify(pubKey string, hash crypto.Hash, msg []byte, sig []byte) error {
	p, _ := pem.Decode([]byte(pubKey))
	if p == nil {
		return errors.New("decode pem failed...")
	}

	key, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return err
	}

	h := hash.New()
	h.Write(msg)
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(key.(*rsa.PublicKey), hash, hashed[:], sig)
}

func GenRSAKeyPair(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	publicKey := privateKey.Public()

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)

	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	return string(pem.EncodeToMemory(privBlock)), string(pem.EncodeToMemory(pubBlock)), nil
}

func main() {
	msg := "hello golang"
	sig, err := RSASign("pkcs8_rsa_private_key.pem", crypto.SHA256, []byte(msg))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%x\n", sig)
		err := RSAVerify("rsa_public_key.pem", crypto.SHA256, []byte(msg), sig)
		fmt.Println(err)
	}
}
