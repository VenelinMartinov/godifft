package godifft_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VenelinMartinov/godifft"
)

func TestDiff(t *testing.T) {
	eq := func(a, b byte) bool {
		return a == b
	}
	input := []byte(`mario`)
	dd := godifft.DiffT(input, []byte(`darius`), godifft.DiffTOptions[byte]{Equals: eq})
	assert.Equal(t, godifft.Remove, dd[0].Change)
	assert.Equal(t, godifft.Insert, dd[1].Change)
	assert.Equal(t, godifft.Keep, dd[2].Change)
	assert.Equal(t, godifft.Keep, dd[3].Change)
	assert.Equal(t, godifft.Keep, dd[4].Change)
	assert.Equal(t, godifft.Remove, dd[5].Change)
	assert.Equal(t, godifft.Insert, dd[6].Change)
	assert.Equal(t, godifft.Insert, dd[7].Change)
}

func TestDiffWithIndices(t *testing.T) {
	eq := func(a, b byte) bool {
		return a == b
	}
	input := []byte(`mario`)
	dd := godifft.DiffT(input, []byte(`darius`), godifft.DiffTOptions[byte]{Equals: eq})
	assert.Equal(t, godifft.Edit[byte]{Change: godifft.Remove, Element: 'm', Index: 0}, dd[0])
	assert.Equal(t, godifft.Edit[byte]{Change: godifft.Insert, Element: 'd', Index: 0}, dd[1])
	assert.Equal(t, godifft.Edit[byte]{Change: godifft.Keep, Element: 'a', Index: 1}, dd[2])
	assert.Equal(t, godifft.Edit[byte]{Change: godifft.Keep, Element: 'r', Index: 2}, dd[3])
	assert.Equal(t, godifft.Edit[byte]{Change: godifft.Keep, Element: 'i', Index: 3}, dd[4])
	assert.Equal(t, godifft.Edit[byte]{Change: godifft.Remove, Element: 'o', Index: 4}, dd[5])
	assert.Equal(t, godifft.Edit[byte]{Change: godifft.Insert, Element: 'u', Index: 4}, dd[6])
	assert.Equal(t, godifft.Edit[byte]{Change: godifft.Insert, Element: 's', Index: 5}, dd[7])
}

func TestDiffMap(t *testing.T) {
	m1 := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	m2 := map[string]string{
		"a": "1",
		"b": "3",
		"d": "2",
	}
	actual := godifft.DiffMapT(&testStringDiffer{}, m1, m2)
	assert.Equal(t, actual, map[string]string{
		"b": "s/2/3",
		"c": "-3",
		"d": "+2",
	})
}

type testStringDiffer struct{}

func (*testStringDiffer) Added(x string) string {
	return "+" + x
}

func (*testStringDiffer) Removed(x string) string {
	return "-" + x
}

func (*testStringDiffer) Diff(x, y string) (string, bool) {
	if x == y {
		return "", false
	}
	return "s/" + x + "/" + y, true
}

var _ godifft.Differ[string, string] = (*testStringDiffer)(nil)

// func TestDiffTree(t *testing.T) {
// 	t1 := map[string]interface{}{
// 		"x": "ok",
// 		"m": map[string]interface{}{
// 			"a": 1,
// 			"b": 2,
// 			"c": 3,
// 		},
// 	}
// 	t2 := map[string]interface{}{
// 		"m": map[string]interface{}{
// 			"a": 1,
// 			"b": 3,
// 			"d": 2,
// 		},
// 	}
// 	actual, neq := godifft.DiffTree[interface{}](&treeDiffer{}, reflect.DeepEqual, t1, t2)
// 	assert.True(t, neq)
// 	assert.Equal(t, map[string]interface{}{
// 		"m": map[string]interface{}{
// 			"b": map[string]interface{}{"+": 3, "-": 2},
// 			"c": map[string]interface{}{"-": 3},
// 			"d": map[string]interface{}{"+": 2},
// 		},
// 		"x": map[string]interface{}{"-": "ok"},
// 	}, actual)
// }

type treeDiffer struct{}

var _ godifft.Differ[interface{}, interface{}] = (*treeDiffer)(nil)

func (*treeDiffer) Added(x interface{}) interface{} {
	return map[string]interface{}{"+": x}
}

func (*treeDiffer) Removed(x interface{}) interface{} {
	return map[string]interface{}{"-": x}
}

func (*treeDiffer) Diff(x, y interface{}) (interface{}, bool) {
	if reflect.DeepEqual(x, y) {
		return nil, false
	}
	return map[string]interface{}{
		"-": x,
		"+": y,
	}, true
}
