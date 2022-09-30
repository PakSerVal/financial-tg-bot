package unknown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Success(t *testing.T) {
	command := New()

	res, err := command.Process("unknown command")
	assert.NoError(t, err)
	assert.Equal(t, "не знаю такую команду", res)
}
