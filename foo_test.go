package go_ci_action

import (
	"testing"
)

func Test_bar(t *testing.T) {

	if v := bar(); v != nil {
		t.Errorf("bar() should return nil, got %v", v)
	}
}
