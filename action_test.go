package brest_test

import (
	"testing"

	"github.com/aptogeo/brest"
	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	assert.NotEqual(t, brest.None, brest.Get)
	assert.NotEqual(t, brest.None, brest.Post)
	assert.NotEqual(t, brest.None, brest.Put)
	assert.NotEqual(t, brest.None, brest.Patch)
	assert.NotEqual(t, brest.None, brest.Delete)
	assert.NotEqual(t, brest.None, brest.All)
	assert.NotEqual(t, brest.Get, brest.Post)
	assert.NotEqual(t, brest.Get, brest.Put)
	assert.NotEqual(t, brest.Get, brest.Patch)
	assert.NotEqual(t, brest.Get, brest.Delete)
	assert.NotEqual(t, brest.Get, brest.All)
	assert.NotEqual(t, brest.Post, brest.Put)
	assert.NotEqual(t, brest.Post, brest.Patch)
	assert.NotEqual(t, brest.Post, brest.Delete)
	assert.NotEqual(t, brest.Post, brest.All)
	assert.NotEqual(t, brest.Put, brest.Patch)
	assert.NotEqual(t, brest.Put, brest.Delete)
	assert.NotEqual(t, brest.Put, brest.All)
	assert.NotEqual(t, brest.Patch, brest.Delete)
	assert.NotEqual(t, brest.Patch, brest.All)
	assert.NotEqual(t, brest.Delete, brest.All)
	assert.Equal(t, brest.All, brest.Get+brest.Post+brest.Put+brest.Patch+brest.Delete)
	assert.Equal(t, 0, int(brest.None))
	assert.Equal(t, 1, int(brest.Get))
	assert.Equal(t, 2, int(brest.Post))
	assert.Equal(t, 4, int(brest.Put))
	assert.Equal(t, 8, int(brest.Patch))
	assert.Equal(t, 16, int(brest.Delete))
	assert.Equal(t, 31, int(brest.All))
}
