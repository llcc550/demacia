package basefunc

import "golang.org/x/crypto/bcrypt"

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	res, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(res)
}
