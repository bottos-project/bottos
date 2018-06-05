package safemath

import (
	"math"
	"testing"
	//"fmt"

	"github.com/stretchr/testify/assert"
)

func TestOverFlow(t *testing.T) {
	var u64 uint64
	var u64max uint64
	var u64zero uint64

	u64max = math.MaxUint64
	u64zero = uint64(0)

	u64 = u64max
	assert.Equal(t, u64max, u64)

	// overflow
	u64 = u64 + 1
	assert.Equal(t, u64zero, u64)

	// overflow
	u64 = 0
	u64 = u64 - 1
	assert.Equal(t, u64max, u64)
}

func TestSafeMath(t *testing.T) {
	var u64 uint64
	var u64max uint64
	var u64zero uint64

	u64max = math.MaxUint64
	u64zero = uint64(0)

	u64 = u64max
	assert.Equal(t, u64max, u64)

	// overflow
	var a, b, c uint64
	a = u64max
	b = uint64(1)
	c, err := Uint64Add(a, b)
	assert.NotNil(t, err)

	// overflow
	u64 = 0
	u64 = u64 - 1
	assert.Equal(t, u64max, u64)
	a = u64zero
	b = uint64(1)
	c, err = Uint64Sub(a, b)
	assert.NotNil(t, err)

	// normal add
	a = 9999999
	b = 88888888
	c, err = Uint64Add(a, b)
	assert.Equal(t, a+b, c)

	// normal sub
	a = 99999999
	b = 88888888
	c, err = Uint64Sub(a, b)
	assert.Equal(t, a-b, c)
}
