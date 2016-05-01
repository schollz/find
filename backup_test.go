package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	assert.Equal(t, dumpFingerprints("testdb"), nil)
}
