package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

// https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes/

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
		go log.Printf("Cannot crate cipher block from key %v.", key)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		go log.Printf("Cannot create new GCM from block %v.", block)
	}
	return gcm
}

func encrypt(toEncrypt []byte, pass string) []byte {
	gcm := prepare(pass)

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		go log.Printf("Cannot create nonce from gcm %v.", nonce)
	}
	return gcm.Seal(nonce, nonce, toEncrypt, nil)
}

func decrypt(encrpyted []byte, pass string) []byte {
	gcm := prepare(pass)
	nonceSize := gcm.NonceSize()
	nonce, ciphered := encrpyted[:nonceSize], encrpyted[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, ciphered, nil)
	if err != nil {
		go log.Printf("Cannot decrypt data.")
	}
	return decrypted
}
