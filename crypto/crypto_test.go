package crypto

import "testing"
import "github.com/stretchr/testify/assert"

func TestHashDataToStringShouldProduceCorrectOutput(t *testing.T) {
	result, err := HashDataToString([]byte("testing 123"))
	assert.NoError(t, err)

	expected := "3XKaGNNw9eqMbJX9ZfCCw8ux7xsy476Kz1PDR6sh9zPc6wAqWZQcM6iLb3LReXbGt4UwtTT6qzSUoSqZbL6oG3mb"
	assert.Equal(t, expected, result)
}
