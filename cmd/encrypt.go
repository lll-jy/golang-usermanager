package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

/*func encrypt(file multipart.File, pass string, info *PageInfo, client *encryption.Client) {
	targetDir := "../../../Desktop/EntryTask/entry-task/test/data/upload" // EXTEND: May set to some cloud space
	tempFile, err := ioutil.TempFile(targetDir, "upload-*.jpeg")
	if err != nil {
		log.Println("Error generating temporary file.")
		log.Println(err)
	}
	defer tempFile.Close()
	img, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Error reading file.")
	}
	// img := utils.ReadFile(file)
	cipherImg := client.EncryptAES(img, []byte(pass))

	// tempFile.Write(fileBytes)
	dirs := strings.Split(tempFile.Name(), "/")
	info.TempUser.PhotoUrl = fmt.Sprintf("test/data/upload/%s", dirs[len(dirs)-1]) // EXTEND: same as above
	err = utils.WriteFile(cipherImg, info.TempUser.PhotoUrl)

	// err = utils.WriteFile(cipherImg, "??")
	if err != nil {
		log.Printf("Cannot encrypt file %s", file)
	}
}

func decrypt(encrypted string, pass string, info *PageInfo, client *encryption.Client) {
	encryptedImg := utils.ReadFile("test/data/upload") // EXTEND:
	tempImg := client.DecryptAES(encryptedImg, []byte(pass))
	err := utils.WriteFile(tempImg, fmt.Sprintf("test/data/temp/temp-%s.jpeg", info.User.Name))
	if err != nil {
		log.Printf("Cannot decrypt file %s", encrypted)
	}
	info.TempUser.PhotoUrl = "test/data/temp/temp-%s.jpeg"
}
*/
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
		log.Printf("Cannot crate cipher block from key %v.", key)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Cannot create new GCM from block %v.", block)
	}
	return gcm
}

func encrypt(toEncrypt []byte, pass string) []byte {
	gcm := prepare(pass)

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Cannot create nonce from gcm %v.", nonce)
	}
	return gcm.Seal(nonce, nonce, toEncrypt, nil)
}

func decrypt(encrpyted []byte, pass string) []byte {
	gcm := prepare(pass)
	nonceSize := gcm.NonceSize()
	nonce, ciphered := encrpyted[:nonceSize], encrpyted[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, ciphered, nil)
	if err != nil {
		log.Printf("Cannot decrypt data.")
	}
	return decrypted
}
