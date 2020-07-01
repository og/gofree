package f

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUID(t *testing.T) {
	assert.Equal(t, 36, len(UUID()))
}