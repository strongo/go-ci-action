package go_ci_action

import (
	// we add the 'assert' package to the imports to have a go.sum file for caching in CI
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_bar(t *testing.T) {
	v := bar()
	assert.Nil(t, v)
}
