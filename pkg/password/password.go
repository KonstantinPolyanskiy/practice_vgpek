package password

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	pb := []byte(password)

	hashed, err := bcrypt.GenerateFromPassword(pb, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
