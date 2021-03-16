package secure

import (
	"os"
	"strconv"

	"github.com/contractGuru/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) ([]byte, error) {
	hashValue, err := strconv.Atoi(os.Getenv("HASH_VALUE"))

	if err != nil {
		logger.Error.Println(err.Error())
	}

	return bcrypt.GenerateFromPassword([]byte(password), hashValue)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
