package group

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {

	{ // Con valores
		g := Int([]int{1, 2, 3})
		v, err := g.Value()
		assert.Nil(t, err)
		assert.Equal(t, []byte("{1,2,3}"), v, string(v.([]byte)))
	}

	{ // Vacio
		g := Int([]int{})
		v, err := g.Value()
		assert.Nil(t, err)
		assert.Equal(t, []byte("{}"), v, string(v.([]byte)))
	}
}
