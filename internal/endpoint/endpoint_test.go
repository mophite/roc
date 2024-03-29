package endpoint

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewEndpoint(t *testing.T) {
    e, err := NewEndpoint("12", "test", "127.0.0.1:8080")
    assert.Nil(t, err)

    assert.Equal(t, "goroc/test/v1.0.0/127.0.0.1:8080", e.Absolute)
    assert.Equal(t, "test/v1.0.0", e.Scope)
}
