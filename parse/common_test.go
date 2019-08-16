package parse

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestGetLastBitPos(t *testing.T) {
	assert.Equal(t, GetLastBitPos(0x82), 2)
}
