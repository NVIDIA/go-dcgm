package dcgm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupHandle(t *testing.T) {
	gh := GroupHandle{}
	assert.Equal(t, uintptr(0), gh.GetHandle(), "value mismatch")

	inputs := []uintptr{1000, 0, 1, 10, 11, 50, 100, 1939902, 9992932938239, 999999999999999999}

	for _, input := range inputs {
		gh.SetHandle(input)
		assert.Equal(t, input, gh.GetHandle(), "values mismatch")
	}
}
