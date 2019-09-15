package crypto

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func IsPasswordMatch(actual string, verify string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(actual), []byte(verify))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
