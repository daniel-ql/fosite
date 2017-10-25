package hmac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomBytes(t *testing.T) {
	bytes, err := RandomBytes(128)
	assert.NoError(t, err)
	assert.Len(t, bytes, 128)
}

func TestPseudoRandomness(t *testing.T) {
	runs := 65536
	results := map[string]bool{}
	for i := 0; i < runs; i++ {
		bytes, err := RandomBytes(128)
		assert.NoError(t, err)

		_, ok := results[string(bytes)]
		assert.False(t, ok)
		results[string(bytes)] = true
	}
}
