# v3 singleflight 

## `/:code`接口压测结果 
为方便测试，将key的过期时间设置为1s

### snowflake-有singleflight-同步incr
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   197.04ms   70.41ms 812.39ms   82.67%
    Req/Sec   127.83     23.92   210.00     72.67%
  15299 requests in 30.04s, 4.23MB read
Requests/sec:    509.30
Transfer/sec:    144.24KB
### snowflake-无singleflight-同步incr
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   202.70ms   76.05ms   1.17s    82.30%
    Req/Sec   124.53     25.73   191.00     72.33%
  14906 requests in 30.05s, 4.12MB read
Requests/sec:    496.08
Transfer/sec:    140.49KB
### snowflake-有singleflight-异步incr
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections

  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.00ms    4.51ms  54.71ms   73.16%
    Req/Sec     3.22k   422.34     4.40k    70.67%
  385453 requests in 30.08s, 106.60MB read
Requests/sec:  12814.38
Transfer/sec:      3.54MB
### snowflake-无singleflight-异步incr
Running 30s test @ http://127.0.0.1:8888
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.67ms    5.38ms  87.37ms   77.64%
    Req/Sec     3.01k   469.05     4.07k    70.98%
  359952 requests in 30.09s, 99.55MB read
Requests/sec:  11962.53
Transfer/sec:      3.31MB