package utils

import "golang.org/x/crypto/bcrypt"

// EncryptPassword - encrypt the user password
func EncryptPassword(pswd string) (string, error) {
	encryptByte, err := bcrypt.GenerateFromPassword([]byte(pswd), 10)
	if err != nil {
		return "", err
	}
	return string(encryptByte), nil
}
