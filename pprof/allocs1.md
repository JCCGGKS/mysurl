heap profile: 4: 1052784 [31: 1122032] @ heap/1048576
0: 0 [1: 576] @ 0x479013 0x409afd 0x409afe 0x409a4f 0x408a6d 0x40d525 0x7957f9 0x794cf8 0x7961d3 0x79604b 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x7957f8	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0xf8	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:212
#	0x794cf7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7961d2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x79604a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [0: 0] @ 0x4555f9 0x46cf33 0x4667fb 0x5e0585 0xa6d55b 0x764b69 0x766a64 0x7856ae 0x763065 0x4836c1
#	0x5e0584	runtime/trace.Start+0x84		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/trace/trace.go:125
#	0xa6d55a	net/http/pprof.Trace+0x2ba		/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/pprof/pprof.go:183
#	0x764b68	net/http.HandlerFunc.ServeHTTP+0x28	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2294
#	0x766a63	net/http.(*ServeMux).ServeHTTP+0x1c3	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2822
#	0x7856ad	net/http.serverHandler.ServeHTTP+0x8d	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:3301
#	0x763064	net/http.(*conn).serve+0x624		/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2102

0: 0 [0: 0] @ 0x794905 0x79574a 0x794cf8 0x7961d3 0x79604b 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x794904	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x2c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:101
#	0x795749	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794cf7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7961d2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x79604a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 32] @ 0x4099f2 0x408a6d 0x40c838 0x899507 0x89a18e 0x898de2 0x898d73 0x4550f8 0x4464e5 0x4463ce 0x4836c1
#	0x899506	github.com/prometheus/client_golang/prometheus.(*Registry).Register+0x346	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/prometheus/client_golang@v1.23.2/prometheus/registry.go:306
#	0x89a18d	github.com/prometheus/client_golang/prometheus.(*Registry).MustRegister+0x4d	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/prometheus/client_golang@v1.23.2/prometheus/registry.go:405
#	0x898de1	github.com/prometheus/client_golang/prometheus.MustRegister+0x81		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/prometheus/client_golang@v1.23.2/prometheus/registry.go:177
#	0x898d72	github.com/prometheus/client_golang/prometheus.init.0+0x12			/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/prometheus/client_golang@v1.23.2/prometheus/registry.go:61
#	0x4550f7	runtime.doInit1+0xd7								/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:7353
#	0x4464e4	runtime.doInit+0x344								/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:7320
#	0x4463cd	runtime.main+0x22d								/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:254

0: 0 [1: 16] @ 0x794825 0x79481e 0x79574a 0x794cf8 0x7961d3 0x79604b 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x794824	bufio.(*Scanner).Text+0x1e4						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:115
#	0x79481d	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1dd		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:90
#	0x795749	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794cf7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7961d2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x79604a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 256] @ 0x794905 0x79574a 0x794cf8 0x7961d3 0x79604b 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x794904	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x2c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:101
#	0x795749	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794cf7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7961d2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x79604a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 256] @ 0x794905 0x79635c 0x79605d 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x794904	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x2c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:101
#	0x79635b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x79605c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 16] @ 0x794825 0x79481e 0x79635c 0x79605d 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x794824	bufio.(*Scanner).Text+0x1e4						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:115
#	0x79481d	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1dd		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:90
#	0x79635b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x79605c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [2: 6144] @ 0x794825 0x79481e 0x79635c 0x79605d 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x794824	bufio.(*Scanner).Text+0x1e4						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:115
#	0x79481d	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1dd		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:90
#	0x79635b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x79605c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 32] @ 0x4099f2 0x408a6d 0x40d525 0x7957f9 0x794cf8 0x7961d3 0x79604b 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x7957f8	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0xf8	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:212
#	0x794cf7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7961d2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x79604a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 128] @ 0x794905 0x79574a 0x794cf8 0x7961d3 0x79604b 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x794904	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x2c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:101
#	0x795749	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794cf7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7961d2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x79604a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [8: 32768] @ 0x541cd9 0x794805 0x79574a 0x794cf8 0x7961d3 0x79604b 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x541cd8	bufio.(*Scanner).Scan+0x378						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:209
#	0x794804	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:89
#	0x795749	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794cf7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7961d2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x79604a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [0: 0] @ 0x4555b6 0x46cf6e 0x4667fb 0x5e0585 0xa6d55b 0x764b69 0x766a64 0x7856ae 0x763065 0x4836c1
#	0x5e0584	runtime/trace.Start+0x84		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/trace/trace.go:125
#	0xa6d55a	net/http/pprof.Trace+0x2ba		/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/pprof/pprof.go:183
#	0x764b68	net/http.HandlerFunc.ServeHTTP+0x28	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2294
#	0x766a63	net/http.(*ServeMux).ServeHTTP+0x1c3	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2822
#	0x7856ad	net/http.serverHandler.ServeHTTP+0x8d	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:3301
#	0x763064	net/http.(*conn).serve+0x624		/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2102

0: 0 [1: 64] @ 0x794825 0x79481e 0x79635c 0x79605d 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x794824	bufio.(*Scanner).Text+0x1e4						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:115
#	0x79481d	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1dd		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:90
#	0x79635b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x79605c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [7: 28672] @ 0x541cd9 0x794805 0x79635c 0x79605d 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x541cd8	bufio.(*Scanner).Scan+0x378						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:209
#	0x794804	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:89
#	0x79635b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x79605c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 288] @ 0x479013 0x4089e5 0x4089d8 0x40d4f9 0x7957f9 0x794cf8 0x7961d3 0x79604b 0x7990af 0x5c3f73 0x798fab 0x4836c1
#	0x7957f8	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0xf8	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:212
#	0x794cf7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7961d2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x79604a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x7990ae	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3f72	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798faa	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

1: 112 [1: 112] @ 0x4121f4 0x5e304b 0x4836c1
#	0x5e304a	github.com/zeromicro/go-zero/core/proc.init.1.func1+0x2a	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/proc/signals.go:24

1: 2048 [1: 2048] @ 0x44a2d1 0x44adf5 0x44b4d9 0x44b9b8 0x44ba32 0x44de1a 0x44e345 0x48168e
#	0x44a2d0	runtime.allocm+0x90		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2236
#	0x44adf4	runtime.newm+0x34		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2772
#	0x44b4d8	runtime.startm+0x158		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2998
#	0x44b9b7	runtime.handoffp+0x357		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3039
#	0x44ba31	runtime.stoplockedm+0x51	/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3161
#	0x44de19	runtime.schedule+0x39		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3999
#	0x44e344	runtime.park_m+0x284		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:4144
#	0x48168d	runtime.mcall+0x4d		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/asm_amd64.s:459

1: 2048 [1: 2048] @ 0x44a2d1 0x44adf5 0x44b4d9 0x47b88c 0x44da5e 0x44deef 0x449a0d 0x449915 0x481605
#	0x44a2d0	runtime.allocm+0x90		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2236
#	0x44adf4	runtime.newm+0x34		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2772
#	0x44b4d8	runtime.startm+0x158		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2998
#	0x47b88b	runtime.wakep+0xeb		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3145
#	0x44da5d	runtime.resetspinning+0x3d	/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3885
#	0x44deee	runtime.schedule+0x10e		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:4038
#	0x449a0c	runtime.mstart1+0xcc		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:1862
#	0x449914	runtime.mstart0+0x74		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:1808
#	0x481604	runtime.mstart+0x4		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/asm_amd64.s:395

1: 1048576 [1: 1048576] @ 0x4555b6 0x46cf33 0x4667fb 0x5e0585 0xa6d55b 0x764b69 0x766a64 0x7856ae 0x763065 0x4836c1
#	0x5e0584	runtime/trace.Start+0x84		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/trace/trace.go:125
#	0xa6d55a	net/http/pprof.Trace+0x2ba		/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/pprof/pprof.go:183
#	0x764b68	net/http.HandlerFunc.ServeHTTP+0x28	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2294
#	0x766a63	net/http.(*ServeMux).ServeHTTP+0x1c3	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2822
#	0x7856ad	net/http.serverHandler.ServeHTTP+0x8d	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:3301
#	0x763064	net/http.(*conn).serve+0x624		/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2102


# runtime.MemStats
# Alloc = 5878296
# TotalAlloc = 25293392
# Sys = 22631688
# Lookups = 0
# Mallocs = 89773
# Frees = 64391
# HeapAlloc = 5878296
# HeapSys = 16089088
# HeapIdle = 8740864
# HeapInuse = 7348224
# HeapReleased = 5595136
# HeapObjects = 25382
# Stack = 688128 / 688128
# MSpan = 161920 / 179520
# MCache = 9664 / 15704
# BuckHashSys = 1449741
# GCSys = 2882064
# OtherSys = 1327443
# NextGC = 8388608
# LastGC = 1779901970292867594
# PauseNs = [530526 93371 128556 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
# PauseEnd = [1779901840545833282 1779901923733465108 1779901970292867594 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
# NumGC = 3
# NumForcedGC = 0
# GCCPUFraction = 5.937424150565279e-06
# DebugGC = false
# MaxRSS = 41267200
