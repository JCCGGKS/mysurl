# v2 有无缓存性能对比

## `/:code`接口压测结果 

### snowflake-无缓存-有incr
running: SHORT_CODE=2snRqpmjqNi wrk -t4 -c100 -d30s -s get_code.lua http://127.0.0.1:8888
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   211.86ms   60.27ms 627.66ms   74.06%
    Req/Sec   118.19     25.86   200.00     66.58%
  14153 requests in 30.06s, 3.91MB read
Requests/sec:    470.82
Transfer/sec:    133.34KB

### snowflake-有缓存-有incr
running: SHORT_CODE=2snRqpmjqNi wrk -t4 -c100 -d30s -s get_code.lua http://127.0.0.1:8888
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   201.78ms   76.60ms 855.02ms   81.28%
    Req/Sec   125.09     25.75   202.00     71.58%
  14973 requests in 30.04s, 4.14MB read
Requests/sec:    498.41
Transfer/sec:    141.15KB

### snowflake-无缓存-无incr
running: SHORT_CODE=2snRqpmjqNi wrk -t4 -c100 -d30s -s get_code.lua http://127.0.0.1:8888
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    14.02ms   10.10ms 116.70ms   73.61%
    Req/Sec     1.92k   353.93     2.81k    68.33%
  230189 requests in 30.08s, 63.66MB read
Requests/sec:   7651.69
Transfer/sec:      2.12MB

### snowflake-有缓存-无incr
running: SHORT_CODE=2snRqpmjqNi wrk -t4 -c100 -d30s -s get_code.lua http://127.0.0.1:8888
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.62ms    6.38ms  89.64ms   74.46%
    Req/Sec     3.13k   662.21     6.29k    65.47%
  374797 requests in 30.10s, 103.66MB read
Requests/sec:  12451.98
Transfer/sec:      3.44MB



### snowflake-有缓存-无布隆过滤器-访问大量不存在的code
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    16.18ms    9.37ms 108.71ms   71.65%
    Req/Sec     1.59k   233.62     4.85k    74.19%
  190664 requests in 30.09s, 45.09MB read
  Non-2xx or 3xx responses: 190664
Requests/sec:   6336.36
Transfer/sec:      1.50MB

### snowflake-有缓存-有布隆过滤器-访问大量不存在code
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     6.16ms    4.05ms  65.47ms   71.10%
    Req/Sec     4.21k   516.68     5.45k    75.50%
  503950 requests in 30.07s, 119.19MB read
  Non-2xx or 3xx responses: 503950
Requests/sec:  16761.55
Transfer/sec:      3.96MB

### snowflake-有缓存-有布隆过滤器-访问存在code
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    15.73ms   43.67ms 627.48ms   97.32%
    Req/Sec     2.62k   455.28     3.73k    78.71%
  305688 requests in 30.20s, 84.54MB read
Requests/sec:  10121.30
Transfer/sec:      2.80MB
`code存在的情况下，bloom会延长请求链路`