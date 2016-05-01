// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// network.go contains structures and functions for creating networks from slices.

package main

import (
	"encoding/json"
	"strconv"
)

func hasNetwork(network map[string]map[string]bool, macs []string) (string, bool) {
	for n := range network {
		for _, val := range macs {
			if network[n][val] {
				return n, true
			}
		}
	}
	return "none", false
}

func buildNetwork(network map[string]map[string]bool, macs []string) map[string]map[string]bool {
	if len(network) == 0 {
		network["0"] = make(map[string]bool)
		for _, val := range macs {
			network["0"][val] = true
		}
		return network
	}
	networkName, inNetworkAlready := hasNetwork(network, macs)
	if inNetworkAlready {
		for _, val := range macs {
			network[networkName][val] = true
		}

	} else {
		// Iterate network to get new name
		curVal := 0
		for n := range network {
			num, _ := strconv.Atoi(n)
			if num > curVal {
				curVal = num
			}
		}
		curVal++
		networkName := strconv.Itoa(curVal)
		network[networkName] = make(map[string]bool)
		for _, val := range macs {
			network[networkName][val] = true
		}
	}
	return network
}

func mergeNetwork(network map[string]map[string]bool) map[string]map[string]bool {
	for {
		n, m, canMerge := hasMerge(network)

		if canMerge {
			for k := range network[m] {
				network[n][k] = true
			}
			delete(network, m)
		} else {
			break
		}
	}
	return network

}

func hasMerge(network map[string]map[string]bool) (string, string, bool) {
	for n := range network {
		for m := range network {
			if m != n {
				for j := range network[n] {
					for k := range network[m] {
						if j == k {
							return n, m, true
						}
					}
				}
			}
		}
	}
	return "none", "none", false
}

func dumpNetwork(network map[string]map[string]bool) []byte {
	jsonByte, _ := json.Marshal(network)
	return jsonByte
}

func loadNetwork(jsonByte []byte) map[string]map[string]bool {
	res2 := make(map[string]map[string]bool)
	json.Unmarshal(jsonByte, &res2)
	return res2
}

// network := make(map[string]map[string]bool)
// macs := []string{"test", "test2", "test3", "test4"}
// fmt.Println(buildNetwork(network, macs))
// macs = []string{"test6", "test5", "test7"}
// fmt.Println(buildNetwork(network, macs))
// macs = []string{"test6", "test9", "test10"}
// fmt.Println(buildNetwork(network, macs))
// fmt.Println(mergeNetwork(network))
// n2 := dumpNetwork(network)
// fmt.Println(string(n2))
// n2a := loadNetwork(n2)
// fmt.Println(n2a)
