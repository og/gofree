package f

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestIgnorePatternAndEmpty(t *testing.T) {
	assert.Equal(t, EqualIgnoreString("success", "all"),Equal("success"))
	assert.Equal(t, EqualIgnoreString("all", "all"), IgnoreFilter())

	assert.Equal(t, EqualIgnoreEmpty("success"),Equal("success"))
	assert.Equal(t, EqualIgnoreEmpty(""), IgnoreFilter())
}
func TestIgnore(t *testing.T) {
	assert.Equal(t, Ignore(Like("success"), "success" == "all"),Like("success"))
	assert.Equal(t, Ignore(Like("success"), "all" == "all"),IgnoreFilter())

	assert.Equal(t, Ignore(Like("success"), "success" == ""),Like("success"))
	assert.Equal(t, Ignore( Like(""), "" == ""),IgnoreFilter())
}
