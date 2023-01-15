package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHelloTom(t *testing.T) {
	output := HelloTom()
	expectOutput := "Tom"
	if output != expectOutput {
		//t.Errorf("Excepted %s do not match actual %s", expectOutput, output)
		assert.Equal(t, expectOutput, output)
	}
}
