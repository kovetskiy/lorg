package lorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDiscarder_ReturnsDiscarderInstance(t *testing.T) {
	test := assert.New(t)

	instance := NewDiscarder()
	test.IsType((*discarder)(nil), instance)
}
