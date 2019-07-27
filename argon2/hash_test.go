package argon2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFromPassword(t *testing.T) {
	hash, err := GenerateFromPassword([]byte("test"))

	assert.Nil(t, err)

	fmt.Printf("%v", hash)
}
