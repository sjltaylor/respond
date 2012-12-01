package crypto

import (
	"code.google.com/p/gorilla/securecookie"
	"crypto/sha512"
	"encoding/base64"
)

func GenerateRandomBase64String(byteCount int) string {
	return base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(byteCount))
}

func GenerateRandomBytes(byteCount int) []byte {
	return securecookie.GenerateRandomKey(byteCount)
}

func Base64Salt(size int) string {
	bytes := GenerateRandomBytes(size)
	return base64.StdEncoding.EncodeToString(bytes)
}

func Base64SaltAndHash(password string, saltByteCount int) (string, string) {

	b64Salt := Base64Salt(saltByteCount)

	b64Hash := Base64HashOfSaltedPassword(password, b64Salt)

	return b64Salt, b64Hash
}

func Base64HashOfSaltedPassword(plainPassword, b64Salt string) string {

	saltedPassword := append([]byte(plainPassword), []byte(b64Salt)...)

	hash := sha512.New()
	if _, err := hash.Write(saltedPassword); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
