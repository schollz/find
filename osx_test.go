package main

import (
  "fmt"
  "io/ioutil"
)

func ExampleOSX1() {
  dat, _ := ioutil.ReadFile("test/osxOutput.txt")
  results := processOutputOSX(string(dat))
  fmt.Println(results)
  // Output:[{22:86:8c:d5:30:d8 -82} {10:86:8c:d5:30:d8 -81} {00:7f:28:8b:0c:1d -74} {74:44:01:35:46:34 -75} {20:e5:2a:16:79:d4 -73} {4e:7a:8a:1d:7e:cc -70} {e6:89:2c:1a:02:e0 -91} {e8:89:2c:1a:02:e0 -92} {3c:7a:8a:1d:7e:cc -70} {50:6a:03:93:89:e3 -91} {08:86:3b:6d:bc:16 -53} {c0:ff:d4:e6:dd:2b -67} {60:a4:4c:29:8b:44 -88} {08:86:3b:6d:bc:18 -66}]
}