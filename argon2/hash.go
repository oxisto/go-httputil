package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrMismatchedHashAndPassword is returned when password and hash do not match
	ErrMismatchedHashAndPassword = errors.New("argon2: encoded is not the hash of the given password")

	// ErrIncompatibleVersion is returned when hash versions are not compatible
	ErrIncompatibleVersion = errors.New("argon2: version is not compatible")

	// ErrInvalidHash is returned when the encoded hash has an invalid format
	ErrInvalidHash = errors.New("argon2: encoded hash has invalid format")
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
func GenerateFromPassword(password []byte) ([]byte, error) {
	return GenerateFromPasswordWithParams(password, DefaultParams)
}

// GenerateFromPasswordWithParams generates a new argon2id encoded hash from the specified password using the parameters specified
func GenerateFromPasswordWithParams(password []byte, params IDParams) ([]byte, error) {
	// generate the salt
	salt, err := GenerateRandomBytes(params.SaltLength)
	if err != nil {
		return nil, err
	}

	return GenerateFromPasswordWithParamsAndSalt(password, params, salt)
}

// GenerateFromPasswordWithParamsAndSalt generates a new argon2id encoded hash from the specified password and salt using the parameters specified
func GenerateFromPasswordWithParamsAndSalt(password []byte, params IDParams, salt []byte) ([]byte, error) {
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

	return []byte(encoded), nil
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

// CompareHashAndPassword compares an argon2id hashed and encoded password with its possible
// plaintext equivalent. Returns nil on success, or an error on failure.
func CompareHashAndPassword(encoded []byte, password []byte) error {
	params, salt, _, err := newFromHash(encoded)
	if err != nil {
		return err
	}

	other, err := GenerateFromPasswordWithParamsAndSalt(password, *params, salt)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(encoded, other) == 1 {
		return nil
	}

	return ErrMismatchedHashAndPassword
}

func newFromHash(encoded []byte) (params *IDParams, salt []byte, hash []byte, err error) {
	var (
		version int
	)

	parts := strings.Split(string(encoded), "$")

	// hash starts with $argonid$
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return nil, nil, nil, ErrInvalidHash
	}

	params = &IDParams{}

	if _, err = fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	if _, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism); err != nil {
		return nil, nil, nil, err
	}

	if salt, err = base64.RawStdEncoding.DecodeString(parts[4]); err != nil {
		return nil, nil, nil, err
	}

	if hash, err = base64.RawStdEncoding.DecodeString(parts[5]); err != nil {
		return nil, nil, nil, err
	}

	params.SaltLength = uint32(len(salt))
	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}
