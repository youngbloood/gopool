
# gopool
goroutine pool in golang. you can expan the goroutine number dynamic.



# mode
```
                           v
                           |
                           | Send(v)
                           |
               ___________\|/___________
              |  1.get chan roundrobin  |
              |      2.chan<-v          |
              |_________________________|
                   /                \
                  /                  \
                 /                    \
                /                      \
 ______________/________________________\______________
| queue       /                          \             |
|(double link/ list)                      \            |
|    ______ /                              \ ______    |
|   | go 1 |           ->          ->->     | go n |   |         
|   |chan 1|               ...              |chan n|   |
|   |______|           <-          <-<-     |______|   |
|    node 1                                  node n    |
|______________________________________________________|
```


node invoke the function to handle `v`

# usage

```
import "github.com/youngbloood/gopool"

func Run(v interface{})error{
     // todo
     // ...
     return nil
}
```
```
pool := gopool.New(5, Run)
pool.StartGo()                   // now goroutine number is 5 in pool
pool.Expand(2)                   // now goroutine number is 7 in pool
pool.Expand(-4)                  // now goroutine number is 3 in pool
pool.Send(v)
```

# benchmark

## without goroutine pool
```
goos: windows
goarch: amd64
pkg: github.com/youngbloood/gopool
cpu: Intel(R) Core(TM) i3-8100 CPU @ 3.60GHz

BenchmarkSendWithout-4           4481758     271.8 ns/op              8 B/op    1 allocs/op
PASS
ok      github.com/youngbloood/gopool   1.581s
```

## benchmark in gopool
```
goos: windows
goarch: amd64
pkg: github.com/youngbloood/gopool
cpu: Intel(R) Core(TM) i3-8100 CPU @ 3.60GHz

BenchmarkSend/go100-4            3687774     325.9 ns/op              8 B/op    1 allocs/op
BenchmarkSend/go200-4            3940887     320.6 ns/op              8 B/op    1 allocs/op
BenchmarkSend/go300-4            3845986     316.8 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go400-4            3507890	 326.3 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go500-4            3593942	 314.0 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go600-4            3785522	 316.1 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go700-4            3910743	 311.6 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go800-4            4060456	 360.1 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go900-4            4047734	 300.2 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1000-4           3809978	 275.8 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1100-4           3934147	 270.4 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1200-4           4722688	 297.1 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1300-4           4085433	 297.2 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1400-4           3856340	 282.8 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1500-4           3787726	 301.2 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1600-4           3692608	 306.7 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1700-4           3778552	 314.9 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1800-4           3880872	 316.3 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go1900-4           4182049	 283.6 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2000-4           3863949	 282.1 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2100-4           3961070	 299.9 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2200-4           4677076	 299.4 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2300-4           3765348	 305.8 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2400-4           3868771	 299.4 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2500-4           3743061	 299.0 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2600-4           3808707	 302.6 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2700-4           3881643	 302.6 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2800-4           3849794	 302.9 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go2900-4           3693442	 305.7 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3000-4           3774939	 305.9 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3100-4           3822352	 319.0 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3200-4           3492675	 311.2 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3300-4           3681730	 318.7 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3400-4           3615405	 323.6 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3500-4           3150464	 341.4 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3600-4           3490915	 335.0 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3700-4           3577533	 313.8 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go3800-4           3251712	 317.0 ns/op             10 B/op	1 allocs/op
BenchmarkSend/go3900-4           3598915	 328.1 ns/op              8 B/op	1 allocs/op
BenchmarkSend/go4000-4           3090613	 345.1 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4100-4           3266596	 337.4 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4200-4           3115744	 322.6 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4300-4           3389952	 320.6 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4400-4           3320646	 355.0 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4500-4           3580236	 335.3 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4600-4           3308036	 331.5 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4700-4           3102055	 383.3 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4800-4           3784656	 326.5 ns/op              9 B/op	1 allocs/op
BenchmarkSend/go4900-4           3380238	 328.8 ns/op              9 B/op	1 allocs/op
BenchmarkSendWithout-4           4651128	 265.7 ns/op              8 B/op	1 allocs/op
```

## benchmarkParallel in gopool
```
goos: windows
goarch: amd64
pkg: github.com/youngbloood/gopool
cpu: Intel(R) Core(TM) i3-8100 CPU @ 3.60GHz

BenchmarkSendParallel/go100-4	     684810       2078 ns/op               0 B/op	     0 allocs/op
BenchmarkSendParallel/go200-4	     590490       2118 ns/op               0 B/op	     0 allocs/op
BenchmarkSendParallel/go300-4	     595902       2102 ns/op               0 B/op	     0 allocs/op
BenchmarkSendParallel/go400-4	     596533       2103 ns/op               0 B/op	     0 allocs/op
BenchmarkSendParallel/go500-4	     626886       2120 ns/op               0 B/op	     0 allocs/op
BenchmarkSendParallel/go600-4	     519132       2113 ns/op               0 B/op	     0 allocs/op
BenchmarkSendParallel/go700-4	     778765       2107 ns/op               0 B/op	     0 allocs/op
BenchmarkSendParallel/go800-4	     684282       2093 ns/op               0 B/op	     0 allocs/op
BenchmarkSendParallel/go900-4	     596480       2133 ns/op               1 B/op	     0 allocs/op
BenchmarkSendParallel/go1000-4	542811       2128 ns/op               1 B/op	     0 allocs/op
BenchmarkSendParallel/go1100-4	564025       2118 ns/op               1 B/op	     0 allocs/op
BenchmarkSendParallel/go1200-4	705931       2121 ns/op               1 B/op	     0 allocs/op
BenchmarkSendParallel/go1300-4	512174       2136 ns/op               2 B/op	     0 allocs/op
BenchmarkSendParallel/go1400-4	562671       2127 ns/op               2 B/op	     0 allocs/op
BenchmarkSendParallel/go1500-4	588344       2127 ns/op               2 B/op	     0 allocs/op
BenchmarkSendParallel/go1600-4	575725       2133 ns/op               2 B/op	     0 allocs/op
BenchmarkSendParallel/go1700-4	568206       2140 ns/op               2 B/op	     0 allocs/op
BenchmarkSendParallel/go1800-4	556369       2141 ns/op               2 B/op	     0 allocs/op
BenchmarkSendParallel/go1900-4	558002       2142 ns/op               2 B/op	     0 allocs/op
BenchmarkSendParallel/go2000-4	559010       2135 ns/op               2 B/op	     0 allocs/op
BenchmarkSendParallel/go2100-4	570990       2056 ns/op               3 B/op	     0 allocs/op
BenchmarkSendParallel/go2200-4	556256       2155 ns/op              16 B/op	     0 allocs/op
BenchmarkSendParallel/go2300-4	557173       2154 ns/op               3 B/op	     0 allocs/op
BenchmarkSendParallel/go2400-4	557887       2139 ns/op               3 B/op	     0 allocs/op
BenchmarkSendParallel/go2500-4	559482       2159 ns/op               3 B/op	     0 allocs/op
BenchmarkSendParallel/go2600-4	557401       2160 ns/op               3 B/op	     0 allocs/op
BenchmarkSendParallel/go2700-4	553144       2178 ns/op               4 B/op	     0 allocs/op
BenchmarkSendParallel/go2800-4	547492       2148 ns/op               4 B/op	     0 allocs/op
BenchmarkSendParallel/go2900-4	547956       2155 ns/op               4 B/op	     0 allocs/op
BenchmarkSendParallel/go3000-4	549854       2162 ns/op               4 B/op	     0 allocs/op
BenchmarkSendParallel/go3100-4	553161       2179 ns/op               4 B/op	     0 allocs/op
BenchmarkSendParallel/go3200-4	550768       2187 ns/op               4 B/op	     0 allocs/op
BenchmarkSendParallel/go3300-4	554307       2159 ns/op               4 B/op	     0 allocs/op
BenchmarkSendParallel/go3400-4	549740       2245 ns/op               5 B/op	     0 allocs/op
BenchmarkSendParallel/go3500-4	576878       2140 ns/op               4 B/op	     0 allocs/op
BenchmarkSendParallel/go3600-4	542407       2148 ns/op               5 B/op	     0 allocs/op
BenchmarkSendParallel/go3700-4	530734       2174 ns/op               5 B/op	     0 allocs/op
BenchmarkSendParallel/go3800-4	545854       2157 ns/op               5 B/op	     0 allocs/op
BenchmarkSendParallel/go3900-4	550875       1991 ns/op               5 B/op	     0 allocs/op
BenchmarkSendParallel/go4000-4	530872       2073 ns/op               6 B/op	     0 allocs/op
BenchmarkSendParallel/go4100-4	549196       2169 ns/op               6 B/op	     0 allocs/op
BenchmarkSendParallel/go4200-4	527707       2190 ns/op               6 B/op	     0 allocs/op
BenchmarkSendParallel/go4300-4	531418       2190 ns/op               6 B/op	     0 allocs/op
BenchmarkSendParallel/go4500-4	527319       2203 ns/op               7 B/op	     0 allocs/op
BenchmarkSendParallel/go4600-4	535797       2186 ns/op               7 B/op	     0 allocs/op
PASS
ok      github.com/youngbloood/gopool
```

# Q&A
Q: how deal with ErrNotIdle?

A: if recieve ErrNotIdle , may be the goroutine pool size is 0, and invoke `.Expand(add int)` to expan the pool.

Q: how deal with ErrTimeOut?

A: then pool get a channel from queue, and send a v(interface{}) into the queue , it may be timeout. If recieve ErrTimeOut, then you can send many times.

Q: how much size should set in goroutine pool?

A: refer the benchmark.


Q: how to work with queue and queue.Pop()?

A: queue is a double link list, Pop() the head and return head-node, and then push the node into queue's tail, then pool can send v into the head-node's channel.


Q: when invoke function `.Expand(reduce int)` to reduce the goroutine in pool, Will the data be lost?

A: Done is a lazy reduce, when invoke it, it will mark reduce number in queue, after the `Pop()`, close(channel) reduce the goroutine number really.