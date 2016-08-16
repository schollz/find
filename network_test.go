package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNetwork(t *testing.T) {
	newNetwork := make(map[string]map[string]bool)
	newNetwork = buildNetwork(newNetwork, []string{"pie", "ice cream"})
	newNetwork = buildNetwork(newNetwork, []string{"pie", "cocolate syrup"})
	newNetwork = buildNetwork(newNetwork, []string{"orange juice", "water"})
	newNetwork = buildNetwork(newNetwork, []string{"water", "coffee"})
	assert.Equal(t, newNetwork["0"], map[string]bool{"ice cream": true, "cocolate syrup": true, "pie": true})
}

func TestHasNetwork(t *testing.T) {
	newNetwork := make(map[string]map[string]bool)
	newNetwork = buildNetwork(newNetwork, []string{"pie", "ice cream"})
	newNetwork = buildNetwork(newNetwork, []string{"pie", "cocolate syrup"})
	newNetwork = buildNetwork(newNetwork, []string{"orange juice", "water"})
	newNetwork = buildNetwork(newNetwork, []string{"water", "coffee"})
	network, _ := hasNetwork(newNetwork, []string{"water"})
	assert.Equal(t, network, "1")
}
