package tests

import (
  "gotest.tools/assert"
	"testing"
  app "app"
)

func Test_Converstion(t *testing.T) {
	// Arrange
	var cid = app.NewBase64ID(0)

	// Act
	n := uint32(1000000)

	for i := uint32(0); i < n; i++ {
		b64 := cid.FromUint32(i)
		reversed, _ := cid.ToUint32(b64)
		// Assert
		assert.Equal(t, i, reversed)
	}
}

func Test_Counter(t *testing.T) {
	// Arrange
	var cid = app.NewBase64ID(0)

	// Act
	n := uint32(4242)
	c := uint32(0)
	for i := uint32(0); i <= n; i++ {
		c = cid.GetNextUint32()
	}

	// Assert
	assert.Equal(t, c, n)
}
