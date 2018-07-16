package markup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptyTag(t *testing.T) {
	r := Tag("foo")
	assert.IsType(t, r, &RenderResult{})
	assert.Equal(t, r.name, "foo")
	assert.Equal(t, r.data, "")
	assert.Nil(t, r.klass)
	assert.Nil(t, r.markups)
	assert.Nil(t, r.children)
}

func TestNewText(t *testing.T) {
	r := Text("foo")
	assert.IsType(t, r, &RenderResult{})
	assert.Equal(t, r.name, "")
	assert.Equal(t, r.data, "foo")
	assert.Nil(t, r.klass)
	assert.Nil(t, r.markups)
	assert.Nil(t, r.children)
}

func TestNewNest1Tag(t *testing.T) {
	r := Tag("foo",
		Tag("bar"),
	)
	assert.IsType(t, r, &RenderResult{})
	assert.Equal(t, r.name, "foo")
	assert.Equal(t, r.data, "")
	assert.Nil(t, r.klass)
	assert.Nil(t, r.markups)
	assert.Equal(t, len(r.children), 1)
	c := r.children[0]
	assert.Equal(t, c.name, "bar")
	assert.Equal(t, c.data, "")
	assert.Nil(t, c.klass)
	assert.Nil(t, c.markups)
	assert.Nil(t, c.children)
}

func TestNewNestEmptyList(t *testing.T) {
	r := Tag("foo",
		List{},
	)
	assert.IsType(t, r, &RenderResult{})
	assert.Equal(t, r.name, "foo")
	assert.Equal(t, r.data, "")
	assert.Nil(t, r.klass)
	assert.Nil(t, r.markups)
	assert.Nil(t, r.children)
}
func TestNewNestList(t *testing.T) {
	r := Tag("foo",
		List{Tag("bar"), Tag("baz")},
	)
	assert.IsType(t, r, &RenderResult{})
	assert.Equal(t, r.name, "foo")
	assert.Equal(t, r.data, "")
	assert.Nil(t, r.klass)
	assert.Nil(t, r.markups)
	assert.Equal(t, len(r.children), 2)
	c := r.children[0]
	assert.Equal(t, c.name, "bar")
	assert.Equal(t, c.data, "")
	assert.Nil(t, c.klass)
	assert.Nil(t, c.markups)
	assert.Nil(t, c.children)
	c = r.children[1]
	assert.Equal(t, c.name, "baz")
	assert.Equal(t, c.data, "")
	assert.Nil(t, c.klass)
	assert.Nil(t, c.markups)
	assert.Nil(t, c.children)
}
