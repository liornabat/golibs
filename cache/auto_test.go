package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAutoCache(t *testing.T) {
	require := require.New(t)
	c := NewAutoCache(5 * time.Second)
	c.Put("a", "a")
	ok := c.Exist("a")
	require.True(ok)
	val := c.Get("a")
	require.Equal("a", val.(string))
	c.Put("a", "b")
	val = c.Get("a")
	require.Equal("b", val.(string))
	c.Delete("a")
	val = c.Get("a")
	require.Nil(val)
	c.Put("a", "b")
	time.Sleep(6 * time.Second)
	val = c.Get("a")
	require.Nil(val)
}
