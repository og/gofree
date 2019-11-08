package f

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUID(t *testing.T) {
	assert.Equal(t, 36, len(UUID()))
}
func TestConvUUID32(t *testing.T) {
	assert.Equal(t, 32, len(GetUUID32(UUID())))
}
func TestConvUUID64(t *testing.T) {
	uuid32 := GetUUID32(UUID())
	assert.Equal(t, 32, len(uuid32))
	assert.Equal(t, 36, len(GetUUID36(uuid32)))
}