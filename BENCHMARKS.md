# Benchmarks 05/07/2016

Current coverage (`go test -cover`): 34.1%.

To test first start process on one terminal `./gofind`, and then in another window `./testing/testdb.sh`. Then you can close `./gofind` and run with `go test -bench=. -test.benchmem`.

### i7-3370 @ 3.4GHz

2,058 fingerprints
```
BenchmarkTrackFingerprintRoute-8             100          15239873 ns/op         1520979 B/op        386 allocs/op
BenchmarkLearnFingerprintRoute-8             100          14811166 ns/op         1523019 B/op        337 allocs/op
BenchmarkPutFingerprintInDatabase-8          100          18666458 ns/op         1498315 B/op        237 allocs/op
BenchmarkGetFingerprintInDatabase-8         1000           2126121 ns/op           53762 B/op         97 allocs/op
BenchmarkLoadFingerprint-8                100000             15466 ns/op            2457 B/op         32 allocs/op
BenchmarkLoadCompressedFingerprint-8       30000             80320 ns/op           44846 B/op         40 allocs/op
BenchmarkDumpFingerprint-8                200000             11307 ns/op            2372 B/op         44 allocs/op
BenchmarkDumpCompressedFingerprint-8         500           3260405 ns/op         1463900 B/op        127 allocs/op
BenchmarkLoadParameters-8                    100          22610022 ns/op         2662599 B/op       3769 allocs/op
BenchmarkGetParameters-8                     100          14269609 ns/op          952125 B/op       5592 allocs/op
BenchmarkCalculatePosteriors1-8            30000             49136 ns/op            2712 B/op         12 allocs/op
BenchmarkOptimizePriors-8                      3         364207766 ns/op        33952237 B/op      57270 allocs/op
BenchmarkOptimizePriorsThreaded-8              5         234406160 ns/op        32756854 B/op      60000 allocs/op
BenchmarkOptimizePriorsThreadedNot-8           5         248316460 ns/op        32395744 B/op      51277 allocs/op
BenchmarkCrossValidation-8                   100          26085247 ns/op         1025146 B/op       4295 allocs/op
BenchmarkCalculatePriors-8                   100          14561229 ns/op          482638 B/op       1405 allocs/op
```

# Benchmarks 05/01/2016

Current coverage (`go test -cover`): 34.1%.

To test first start process on one terminal `./gofind`, and then in another window `./testing/testdb.sh`. Then you can close `./gofind` and run with `go test -bench=. -test.benchmem`.

### i7-3370 @ 3.4GHz

2,058 fingerprints
```
BenchmarkPutFingerprintInDatabase-8          100          18653644 ns/op         1498396 B/op        237 allocs/op
BenchmarkGetFingerprintInDatabase-8          500           3347812 ns/op           52653 B/op         91 allocs/op
BenchmarkLoadFingerprint-8                100000             12204 ns/op            2457 B/op         32 allocs/op
BenchmarkLoadCompressedFingerprint-8       30000             55108 ns/op           44846 B/op         40 allocs/op
BenchmarkDumpFingerprint-8                200000              8032 ns/op            2372 B/op         44 allocs/op
BenchmarkDumpCompressedFingerprint-8        1000           1400471 ns/op         1463938 B/op        128 allocs/op
BenchmarkLoadParameters-8                    100          15343887 ns/op         2661374 B/op       3771 allocs/op
BenchmarkGetParameters-8                     100          13094985 ns/op          733550 B/op       4273 allocs/op
BenchmarkCalculatePosteriors1-8            30000             48986 ns/op            2711 B/op         12 allocs/op
BenchmarkOptimizePriors-8                     10         169952320 ns/op        26243536 B/op      45237 allocs/op
BenchmarkOptimizePriorsThreaded-8              5         215925640 ns/op        27777283 B/op     103140 allocs/op
BenchmarkOptimizePriorsThreadedNot-8          10         183357620 ns/op        25007952 B/op      40562 allocs/op
BenchmarkCrossValidation-8                   100          16531952 ns/op          819349 B/op       3355 allocs/op
BenchmarkCalculatePriors-8                   200           9649543 ns/op          371202 B/op       1109 allocs/op
```




# Benchmarks 04/16/2016

Current coverage (`go test -cover`): 21.2%.

To test first start process on one terminal `./gofind`, and then in another window `./testing/testdb.sh`. Then you can close `./gofind` and run with `go test -bench=. -test.benchmem`.

### i7-3370 @ 3.4GHz

2,058 fingerprints
```
BenchmarkPutFingerprintInDatabase-8          100          19446095 ns/op         1496173 B/op        212 allocs/op
BenchmarkGetFingerprintInDatabase-8         1000           1544347 ns/op           53459 B/op         94 allocs/op
BenchmarkLoadFingerprint-8                200000             10224 ns/op            2481 B/op         32 allocs/op
BenchmarkLoadCompressedFingerprint-8       50000             48337 ns/op           44866 B/op         40 allocs/op
BenchmarkDumpFingerprint-8                200000              6906 ns/op            2388 B/op         44 allocs/op
BenchmarkDumpCompressedFingerprint-8        2000            894943 ns/op         1463963 B/op        119 allocs/op
BenchmarkLoadParameters-8                    100          13045162 ns/op         1616906 B/op       3752 allocs/op
BenchmarkGetParameters-8                      50          34637830 ns/op         5575188 B/op      33086 allocs/op
BenchmarkCalculatePosteriors1-8            50000             31700 ns/op            2687 B/op         12 allocs/op
BenchmarkOptimizePriors-8                      2         863536600 ns/op        180279556 B/op    335850 allocs/op
BenchmarkOptimizePriorsThreaded-8              3         424329900 ns/op        159378608 B/op    326170 allocs/op
BenchmarkOptimizePriorsThreadedNot-8           3         454433666 ns/op        155175656 B/op    233557 allocs/op
BenchmarkCrossValidation-8                    10         102343250 ns/op         6530782 B/op      26685 allocs/op
BenchmarkCalculatePriors-8                   100          15477659 ns/op          954329 B/op       4706 allocs/op
```



### [Python](https://github.com/schollz/find/tree/python3) vs. GO

Both benchmarked using `testing/testdb.sh` which has 344 fingerprints in 3 locations using Intel i7-3370.

| Version | Fingerprints sent to /learn | Optimizing priors through /calculate |
|---------|-----------------------------|--------------------------------------|
| [Python](https://github.com/schollz/find/tree/python3)  | 15 fingerprints/sec         | 3 calculations/min                   |
| Go      | 76 fingerprints/sec         | 619 calculations/min                 |
