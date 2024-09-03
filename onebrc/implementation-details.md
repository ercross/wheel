## Commit History links (help)
1. [Initial commit](https://github.com/ercross/wheel/commit/60585fc4ec724f91b01059c878e95c1e5af7da7d)
2. [Initial test file generator implementation](https://github.com/ercross/wheel/commit/0ab8c453cb8ff55b04553bd549db498025daf083)
3. [Optimized cpu-time for test file generator by 40%](https://github.com/ercross/wheel/commit/6e4d622c8ef169dc988fc3575ec201a1089f1948)
4. [Row reader initial implementation](https://github.com/ercross/wheel/commit/47898660878d8355bcdb1880362d0ff2c70360a8)

## Implementation Details

- ### Memory requirement analysis for storing 1 billion row at runtime
    if all data from the 1 billion row is loaded into a map[string][]float32 at runtime,
    it is expected to grow up to an approximate size of 40GiB following this calculation below  
    
    ```
    string storage: up to 108bytes (i.e., string-header:8bytes data:100bytes)  
    float32 storage: data consumes 4bytes  
    []float32 slice header: 24bytes (i.e., 8bytes each for pointers to data, length, and capacity)  
    key memory: 108bytes * 10,000 unique station names = 1,080,000bytes (1.08MB)  
    An entry of value consumes: 1 million of (float32) + 24bytes(slice-header) = (1,000,000 * 4) + 24 = 4,000,024 bytes  
    Total memory consumption for 10,000 unique station names: 4,000,024 * 10,000 = 40,000,240,000 bytes = 40GiB  
    Max total memory consumption: 1.08MB + 40GiB = 40.00108 GB  
    ```
  
- ### Concurrent batch processing
    To avoid loading all 1 billion row in memory, measurements are read and processed in chunk.
    4 goroutines (excluding main goroutine) are involved in the process:  
    A goroutine reads measurement lines from file, another goroutine collects measurements into a batches  
    1 goroutine run statistical calculator to obtain (min, mean, max) on each batch  
    1 goroutine merge new batch together with existing batches


## Benchmark (100,000 rows)
`go test -bench=BenchmarkReadMeasurements -run=xxx -cpuprofile cpu.prof`  
goos: darwin  
goarch: amd64  
pkg: onebrc  
cpu: Intel(R) Core(TM) i7-1060NG7 CPU @ 1.20GHz  
BenchmarkReadMeasurements-8           31          38163580 ns/op  
PASS  
ok      onebrc  4.103s  
