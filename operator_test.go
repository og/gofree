package f

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIgnorePatternAndEmpty(t *testing.T) {
	assert.Equal(t, IgnorePattern("success", "all"),Eql("success"))
	assert.Equal(t, IgnorePattern("all", "all"), ignoreFilter())

	assert.Equal(t, IgnoreEmpty("success"),Eql("success"))
	assert.Equal(t, IgnoreEmpty(""), ignoreFilter())
}
func TestIgnore(t *testing.T) {
	assert.Equal(t, Ignore("success" == "all", Like("success")),Like("success"))
	assert.Equal(t, Ignore("all" == "all", Like("success")),ignoreFilter())

	assert.Equal(t, Ignore("success" == "", Like("success")),Like("success"))
	assert.Equal(t, Ignore("" == "", Like("")),ignoreFilter())
}
