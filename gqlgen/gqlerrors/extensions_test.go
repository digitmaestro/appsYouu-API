package gqlerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExtensions(t *testing.T) {
	records := NewExtensions("MyType", "MyCode")

	assert.Equal(t, records["type"], "MyType")
	assert.Equal(t, records["code"], "MyCode")
}

func TestConvertStructToMap(t *testing.T) {
	upper := convertStructToMap(struct {
		TestUpper string
	}{
		TestUpper: "upper",
	})

	assert.Equal(t, upper["TestUpper"], "upper")

	lower := convertStructToMap(struct {
		TestLower string `json:"test_lower"`
	}{
		TestLower: "lower",
	})

	assert.Equal(t, lower["test_lower"], "lower")
}
