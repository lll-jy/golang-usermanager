package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

// https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes/

// prepare prepares the GCM for the password given
func prepare(pass string) cipher.AEAD {
	key := make([]byte, 32)
	length := len(pass)
	for i := 0; i < 32; i++ {
		if i < length {
			key[i] = pass[i]
		} else {
			key[i] = 0
		}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Cannot crate cipher block from key %v: %s", key, err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Cannot create new GCM from block %v: %s", block, err.Error())
	}
	return gcm
}

// encrypt encrypts bytes given with password as key to resulting bytes
func encrypt(toEncrypt []byte, pass string) []byte {
	gcm := prepare(pass)

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Cannot create nonce from gcm %v: %s", nonce, err.Error())
	}
	return gcm.Seal(nonce, nonce, toEncrypt, nil)
}

// decrypt decrypts bytes given with password as key to original bytes
func decrypt(encrypted []byte, pass string) []byte {
	gcm := prepare(pass)
	nonceSize := gcm.NonceSize()
	nonce, ciphered := encrypted[:nonceSize], encrypted[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, ciphered, nil)
	if err != nil {
		log.Printf("Cannot decrypt data: %s", err.Error())
	}
	return decrypted
}
