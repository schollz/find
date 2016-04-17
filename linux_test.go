package main

import (
	"fmt"
	"io/ioutil"
)

func ExampleLinux1() {
	dat, _ := ioutil.ReadFile("test/linuxOutput.txt")
	results := processOutputLinux(string(dat))
	fmt.Println(results)
	// Output: [{80:37:73:ba:f7:dc -25} {80:37:73:87:46:82 -59} {98:6b:3d:d7:84:e0 -60} {08:95:2a:b1:e9:55 -76} {2c:b0:5d:36:e3:b8 -54} {58:20:b1:21:63:9f -62} {70:73:cb:bd:9f:b5 -78} {b8:3e:59:78:35:99 -75} {a0:63:91:2b:9e:64 -59} {e0:3f:49:03:fd:38 -61} {30:8d:99:71:95:c5 -78} {80:37:73:ba:f7:d8 -37} {a0:63:91:2b:9e:65 -55} {80:37:73:87:56:36 -52} {00:1a:1e:46:cd:11 -59} {00:23:69:d4:47:9f -59} {54:65:de:6f:7e:d5 -70} {70:73:cb:bd:9f:b6 -77} {f8:35:dd:0a:da:be -82} {00:1a:1e:46:cd:10 -58} {d4:05:98:57:b3:15 -79}]
}

func ExamplePi1() {
	dat, _ := ioutil.ReadFile("test/pi3Output.txt")
	results := processOutputLinux(string(dat))
	fmt.Println(results)
	// Output: [{70:73:cb:bd:9f:b5 -72} {4c:60:de:fe:e5:24 -80} {80:37:73:ba:f7:d8 -16} {a0:63:91:2b:9e:65 -43} {00:23:69:d4:47:9f -81} {80:37:73:87:56:36 -68} {2c:b0:5d:36:e3:b8 -75} {58:20:b1:21:63:9f -75} {30:8d:99:71:95:c5 -81} {c8:b3:73:25:22:51 -85} {00:1a:1e:46:cd:10 -76} {e0:46:9a:6d:02:ea -91} {08:95:2a:b1:e9:55 -81} {00:1d:d4:7c:bd:30 -91} {8c:09:f4:d3:84:50 -90} {00:1a:1e:46:cd:11 -76}]
}
