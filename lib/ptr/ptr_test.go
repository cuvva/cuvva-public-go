package ptr_test

import (
	"testing"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	v := 42
	p := ptr.Ptr(v)
	assert.NotNil(t, p)
	assert.Equal(t, v, *p)
	assert.NotSame(t, &v, p)

	s := "hello"
	sp := ptr.Ptr(s)
	assert.NotNil(t, sp)
	assert.Equal(t, s, *sp)
}

func TestString(t *testing.T) {
	v := "hello"
	p := ptr.String(v)
	assert.NotNil(t, p)
	assert.Equal(t, v, *p)
}

func TestBool(t *testing.T) {
	for _, v := range []bool{true, false} {
		p := ptr.Bool(v)
		assert.NotNil(t, p)
		assert.Equal(t, v, *p)
	}
}

func TestInt(t *testing.T) {
	v := 7
	p := ptr.Int(v)
	assert.NotNil(t, p)
	assert.Equal(t, v, *p)
}

func TestInt64(t *testing.T) {
	v := int64(1234567890)
	p := ptr.Int64(v)
	assert.NotNil(t, p)
	assert.Equal(t, v, *p)
}

func TestFloat64(t *testing.T) {
	v := 3.14
	p := ptr.Float64(v)
	assert.NotNil(t, p)
	assert.Equal(t, v, *p)
}

func TestTime(t *testing.T) {
	v := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	p := ptr.Time(v)
	assert.NotNil(t, p)
	assert.True(t, p.Equal(v))
}

func TestEqual(t *testing.T) {
	t.Run("int pointers", func(t *testing.T) {
		a := ptr.Int(1)
		b := ptr.Int(1)
		c := ptr.Int(2)

		assert.True(t, ptr.Equal(a, b))
		assert.True(t, ptr.Equal[int](nil, nil))
		assert.False(t, ptr.Equal(a, c))
		assert.False(t, ptr.Equal(a, nil))
		assert.False(t, ptr.Equal(nil, b))
	})

	t.Run("string pointers", func(t *testing.T) {
		a := ptr.String("hello")
		b := ptr.String("hello")
		c := ptr.String("world")

		assert.True(t, ptr.Equal(a, b))
		assert.True(t, ptr.Equal[int](nil, nil))
		assert.False(t, ptr.Equal(a, c))
		assert.False(t, ptr.Equal(a, nil))
		assert.False(t, ptr.Equal(nil, b))
	})

	t.Run("struct pointers", func(t *testing.T) {
		type person struct {
			name string
		}
		a := ptr.Ptr(person{name: "John"})
		b := ptr.Ptr(person{name: "John"})
		c := ptr.Ptr(person{name: "Yoko"})

		assert.True(t, ptr.Equal(a, b))
		assert.False(t, ptr.Equal(a, c))
		assert.True(t, ptr.Equal[int](nil, nil))
		assert.False(t, ptr.Equal(a, nil))
		assert.False(t, ptr.Equal(nil, b))
	})
}

func TestCopy(t *testing.T) {
	t.Run("nil returns nil", func(t *testing.T) {
		assert.Nil(t, ptr.Copy[int](nil))
	})

	t.Run("copies the value into a new pointer", func(t *testing.T) {
		v := 42
		p := &v
		c := ptr.Copy(p)

		assert.NotNil(t, c)
		assert.Equal(t, *p, *c)
		assert.NotSame(t, p, c)
	})

	t.Run("mutating the copy does not affect the original", func(t *testing.T) {
		v := "hello"
		p := &v
		c := ptr.Copy(p)

		*c = "world"

		assert.Equal(t, "hello", *p)
		assert.Equal(t, "world", *c)
	})

	t.Run("struct values are duplicated", func(t *testing.T) {
		type person struct {
			name string
		}
		p := ptr.Ptr(person{name: "John"})
		c := ptr.Copy(p)

		c.name = "Yoko"

		assert.Equal(t, "John", p.name)
		assert.Equal(t, "Yoko", c.name)
	})
}
