
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
