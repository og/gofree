package f

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIgnoreWhenEqual(t *testing.T) {
	assert.Equal(t, IgnoreWhenEqual("success", "all"),Eql("success"))
	assert.Equal(t, IgnoreWhenEqual("all", "all"), Ignore())

	assert.Equal(t, IgnoreWhenEqual("success", ""),Eql("success"))
	assert.Equal(t, IgnoreWhenEqual("", ""), Ignore())
}
func TestIgnoreWhenEqualCustomCallFilter(t *testing.T) {
	assert.Equal(t, IgnoreWhenEqualCustomCallFilter("success", "all", Like),Like("success"))
	assert.Equal(t, IgnoreWhenEqualCustomCallFilter("all", "all", Like),Ignore())

	assert.Equal(t, IgnoreWhenEqualCustomCallFilter("success", "",Like),Like("success"))
	assert.Equal(t, IgnoreWhenEqualCustomCallFilter("", "", Like), Ignore())
}

func TestIgnoreWhenEqualCustomValue(t *testing.T) {
	assert.Equal(t, IgnoreWhenEqualCustomValue("success", "all", Like("success")),Like("success"))
	assert.Equal(t, IgnoreWhenEqualCustomValue("all", "all", Like("success")),Ignore())

	assert.Equal(t, IgnoreWhenEqualCustomValue("success", "",Like("success")),Like("success"))
	assert.Equal(t, IgnoreWhenEqualCustomValue("", "", Like("success")), Ignore())
}
