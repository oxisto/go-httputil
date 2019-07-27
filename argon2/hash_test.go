package argon2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFromPassword(t *testing.T) {
	hash, err := GenerateFromPassword([]byte("test"))

	assert.Nil(t, err)

	t.Logf("%v", hash)
}

func TestCompareHashAndPassword(t *testing.T) {
	encoded := []byte("$argon2id$v=19$m=65536,t=3,p=2$ekLs6AdWsbtp81gzwAzGIg$4xL21LVn8cL3QrJsfYCytV6DwKtbWvMuMzhxcx2oN9U")
	password := []byte("test")
	notPassword := []byte("notPassword")

	params, salt, hash, err := newFromHash(encoded)

	assert.Nil(t, err)
	assert.True(t, len(salt) == 16)

	assert.Equal(t, uint32(65536), params.Memory)
	assert.Equal(t, uint32(3), params.Iterations)
	assert.Equal(t, uint8(2), params.Parallelism)

	assert.Equal(t, uint32(16), params.SaltLength)
	assert.Equal(t, uint32(32), params.KeyLength)

	t.Logf("Encoded: %v Params: %v Salt: %v Hash: %v", encoded, params, salt, hash)

	err = CompareHashAndPassword(encoded, password)

	assert.Nil(t, err)

	err = CompareHashAndPassword(encoded, notPassword)

	assert.Equal(t, ErrMismatchedHashAndPassword, err)
}
