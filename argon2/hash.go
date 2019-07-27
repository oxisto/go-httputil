package argon2

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// IDParams holds all necessary parameters to configure argon2id
type IDParams struct {
	SaltLength  uint32
	Iterations  uint32
	Memory      uint32
	Parallelism uint8
	KeyLength   uint32
}

// DefaultParams holds resonable default parameters for argon2id
var DefaultParams IDParams

func init() {
	DefaultParams = IDParams{
		SaltLength:  16,
		KeyLength:   32,
		Iterations:  3,
		Memory:      64 * 1024,
		Parallelism: 2,
	}
}

// GenerateFromPassword generates a new argon2id encoded hash from the specified password
func GenerateFromPassword(password []byte) (string, error) {
	return GenerateFromPasswordWithParams(password, DefaultParams)
}

// GenerateFromPasswordWithParams generates a new argon2id encoded hash from the specified password using the parameters specified
func GenerateFromPasswordWithParams(password []byte, params IDParams) (string, error) {
	// generate the salt
	salt, err := GenerateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	// calculate the hash
	hash := argon2.IDKey(password, salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength)

	// encode both hash and salt
	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64salt,
		b64hash,
	)

	return encoded, nil
}

// GenerateRandomBytes generates random bytes of the specified size
func GenerateRandomBytes(size uint32) ([]byte, error) {
	salt := make([]byte, size)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}
