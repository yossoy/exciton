package markup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptyTag(t *testing.T) {
	r := Tag("foo")
	tr, ok := r.(*tagRenderResult)
	assert.True(t, ok)
	assert.Equal(t, tr.data, "")
	assert.Nil(t, tr.markups)
	assert.Nil(t, tr.children)
}

func TestNewText(t *testing.T) {
	r := Text("foo")
	tr, ok := r.(*textRenderResult)
	assert.True(t, ok)
	assert.Equal(t, tr.text, "foo")
}

func TestNewNest1Tag(t *testing.T) {
	r := Tag("foo",
		Tag("bar"),
	)
	tr, ok := r.(*tagRenderResult)
	assert.True(t, ok)
	assert.Equal(t, tr.name, "foo")
	assert.Equal(t, tr.data, "")
	assert.Nil(t, tr.markups)
	assert.Equal(t, len(tr.children), 1)
	c := tr.children[0]
	tcr, ok := c.(*tagRenderResult)
	assert.True(t, ok)
	assert.Equal(t, tcr.name, "bar")
	assert.Equal(t, tcr.data, "")
	assert.Nil(t, tcr.markups)
	assert.Nil(t, tcr.children)
}

func TestNewNestEmptyList(t *testing.T) {
	r := Tag("foo",
		List{},
	)
	tr, ok := r.(*tagRenderResult)
	assert.True(t, ok)
	assert.Equal(t, tr.name, "foo")
	assert.Equal(t, tr.data, "")
	assert.Nil(t, tr.markups)
	assert.Nil(t, tr.children)
}
func TestNewNestList(t *testing.T) {
	r := Tag("foo",
		List{Tag("bar"), Tag("baz")},
	)
	tr, ok := r.(*tagRenderResult)
	assert.True(t, ok)
	assert.Equal(t, tr.name, "foo")
	assert.Equal(t, tr.data, "")
	assert.Nil(t, tr.markups)
	assert.Equal(t, len(tr.children), 2)
	c := tr.children[0]
	tcr, ok := c.(*tagRenderResult)
	assert.True(t, ok)
	assert.Equal(t, tcr.name, "bar")
	assert.Equal(t, tcr.data, "")
	assert.Nil(t, tcr.markups)
	assert.Nil(t, tcr.children)
	c = tr.children[1]
	tcr, ok = c.(*tagRenderResult)
	assert.True(t, ok)
	assert.Equal(t, tcr.name, "baz")
	assert.Equal(t, tcr.data, "")
	assert.Nil(t, tcr.markups)
	assert.Nil(t, tcr.children)
}
