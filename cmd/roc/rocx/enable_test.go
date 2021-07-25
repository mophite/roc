package rocx

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGetLatestTag(t *testing.T) {
    tag := getLatestTag()
    assert.NotEmpty(t, tag)
    t.Log(tag)
}
