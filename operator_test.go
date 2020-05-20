package f

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestIgnorePatternAndEmpty(t *testing.T) {
	assert.Equal(t, IgnorePattern(Eql, "success", "all"),Eql("success"))
	assert.Equal(t, IgnorePattern(Eql, "all", "all"), IgnoreFilter())

	assert.Equal(t, IgnoreEmpty(Eql, "success"),Eql("success"))
	assert.Equal(t, IgnoreEmpty(Eql, ""), IgnoreFilter())
}
func TestIgnore(t *testing.T) {
	assert.Equal(t, Ignore(Like("success"), "success" == "all"),Like("success"))
	assert.Equal(t, Ignore(Like("success"), "all" == "all"),IgnoreFilter())

	assert.Equal(t, Ignore(Like("success"), "success" == ""),Like("success"))
	assert.Equal(t, Ignore( Like(""), "" == ""),IgnoreFilter())
}
