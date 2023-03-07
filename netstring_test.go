package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNetstring(t *testing.T) {
	r := strings.NewReader("1:a,5:bbbbb,11:aaaaaaaaaaa,")

	data, err := ParseNetstring(r)
	assert.NoError(t, err)
	assert.Equal(t, []byte("a"), data)

	data, err = ParseNetstring(r)
	assert.NoError(t, err)
	assert.Equal(t, []byte("bbbbb"), data)

	data, err = ParseNetstring(r)
	assert.NoError(t, err)
	assert.Equal(t, []byte("aaaaaaaaaaa"), data)

	data, err = ParseNetstring(r)
	assert.NoError(t, err)
	assert.Nil(t, data)
}
