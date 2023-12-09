package difft_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/t0yv0/godifft"
)

func TestDiff(t *testing.T) {
	eq := func(a, b byte) bool {
		return a == b
	}
	input := []byte(`mario`)
	dd := difft.DiffT(input, []byte(`darius`), difft.DiffTOptions[byte]{Equals: eq})
	assert.Equal(t, difft.Remove, dd[0].Change)
	assert.Equal(t, difft.Insert, dd[1].Change)
	assert.Equal(t, difft.Keep, dd[2].Change)
	assert.Equal(t, difft.Keep, dd[3].Change)
	assert.Equal(t, difft.Keep, dd[4].Change)
	assert.Equal(t, difft.Remove, dd[5].Change)
	assert.Equal(t, difft.Insert, dd[6].Change)
	assert.Equal(t, difft.Insert, dd[7].Change)
}
