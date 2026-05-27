heap profile: 5: 14720 [67: 201832] @ heap/1048576
0: 0 [1: 24] @ 0x4d2705 0x4d44f6 0x4d44fe 0x50268b 0x502674 0x503474 0x503437 0x503406 0x501dfe 0x7945e9 0x7945d0 0x79562a 0x794bd8 0x7960b3 0x795f2b 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x4d2704	syscall.ByteSliceFromString+0x84					/home/fanqicheng/.gvm/gos/go1.24.4/src/syscall/syscall.go:52
#	0x4d44f5	syscall.BytePtrFromString+0x35						/home/fanqicheng/.gvm/gos/go1.24.4/src/syscall/syscall.go:68
#	0x4d44fd	syscall.openat+0x3d							/home/fanqicheng/.gvm/gos/go1.24.4/src/syscall/zsyscall_linux_amd64.go:94
#	0x50268a	syscall.Open+0x2a							/home/fanqicheng/.gvm/gos/go1.24.4/src/syscall/syscall_linux.go:284
#	0x502673	os.open+0x13								/home/fanqicheng/.gvm/gos/go1.24.4/src/os/file_open_unix.go:15
#	0x503473	os.openFileNolog.func1+0x93						/home/fanqicheng/.gvm/gos/go1.24.4/src/os/file_unix.go:279
#	0x503436	os.ignoringEINTR+0x56							/home/fanqicheng/.gvm/gos/go1.24.4/src/os/file_posix.go:251
#	0x503405	os.openFileNolog+0x25							/home/fanqicheng/.gvm/gos/go1.24.4/src/os/file_unix.go:278
#	0x501dfd	os.OpenFile+0x3d							/home/fanqicheng/.gvm/gos/go1.24.4/src/os/file.go:392
#	0x7945e8	os.Open+0xc8								/home/fanqicheng/.gvm/gos/go1.24.4/src/os/file.go:370
#	0x7945cf	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0xaf		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:81
#	0x795629	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794bd7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7960b2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x795f2a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [21: 86016] @ 0x541bb9 0x7946e5 0x79623c 0x795f3d 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x541bb8	bufio.(*Scanner).Scan+0x378						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:209
#	0x7946e4	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:89
#	0x79623b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x795f3c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [2: 128] @ 0x794705 0x7946fe 0x79623c 0x795f3d 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x794704	bufio.(*Scanner).Text+0x1e4						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:115
#	0x7946fd	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1dd		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:90
#	0x79623b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x795f3c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 32] @ 0x51e4b5 0x7956a9 0x794bd8 0x7960b3 0x795f2b 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x51e4b4	strings.Fields+0x74							/home/fanqicheng/.gvm/gos/go1.24.4/src/strings/strings.go:402
#	0x7956a8	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0xc8	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:207
#	0x794bd7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7960b2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x795f2a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 288] @ 0x478f73 0x4089e5 0x4089d8 0x40d4f9 0x7956d9 0x794bd8 0x7960b3 0x795f2b 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x7956d8	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0xf8	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:212
#	0x794bd7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7960b2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x795f2a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 48] @ 0x47901a 0x47902f 0x79563e 0x794bd8 0x7960b3 0x795f2b 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x79563d	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x5d	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:205
#	0x794bd7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7960b2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x795f2a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 16] @ 0x7947e5 0x79562a 0x794bd8 0x7960b3 0x795f2b 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x7947e4	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x2c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:101
#	0x795629	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794bd7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7960b2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x795f2a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 176] @ 0x51e4b5 0x796285 0x795f3d 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x51e4b4	strings.Fields+0x74							/home/fanqicheng/.gvm/gos/go1.24.4/src/strings/strings.go:402
#	0x796284	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x84	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:138
#	0x795f3c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 64] @ 0x7947e5 0x79623c 0x795f3d 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x7947e4	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x2c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:101
#	0x79623b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x795f3c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [0: 0] @ 0x794590 0x79562a 0x794bd8 0x7960b3 0x795f2b 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x79458f	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x6f		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:76
#	0x795629	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794bd7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7960b2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x795f2a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [3: 1728] @ 0x478f73 0x409afd 0x409afe 0x409a4f 0x408a6d 0x40d525 0x7956d9 0x794bd8 0x7960b3 0x795f2b 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x7956d8	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0xf8	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:212
#	0x794bd7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7960b2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x795f2a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 256] @ 0x7947e5 0x79623c 0x795f3d 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x7947e4	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x2c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:101
#	0x79623b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x795f3c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [12: 36864] @ 0x794705 0x7946fe 0x79623c 0x795f3d 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x794704	bufio.(*Scanner).Text+0x1e4						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:115
#	0x7946fd	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1dd		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:90
#	0x79623b	github.com/zeromicro/go-zero/core/stat/internal.systemCpuUsage+0x3b	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:132
#	0x795f3c	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x3c		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:43
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [15: 61440] @ 0x541bb9 0x7946e5 0x79562a 0x794bd8 0x7960b3 0x795f2b 0x798f8f 0x5c3e53 0x798e8b 0x4835a1
#	0x541bb8	bufio.(*Scanner).Scan+0x378						/home/fanqicheng/.gvm/gos/go1.24.4/src/bufio/scan.go:209
#	0x7946e4	github.com/zeromicro/go-zero/core/iox.ReadTextLines+0x1c4		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/iox/read.go:89
#	0x795629	github.com/zeromicro/go-zero/core/stat/internal.currentCgroupV2+0x49	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:200
#	0x794bd7	github.com/zeromicro/go-zero/core/stat/internal.currentCgroup+0x17	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cgroup_linux.go:42
#	0x7960b2	github.com/zeromicro/go-zero/core/stat/internal.cpuUsage+0x12		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:73
#	0x795f2a	github.com/zeromicro/go-zero/core/stat/internal.RefreshCpu+0x2a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/internal/cpu_linux.go:38
#	0x798f8e	github.com/zeromicro/go-zero/core/stat.init.1.func1.1+0xe		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:35
#	0x5c3e52	github.com/zeromicro/go-zero/core/threading.RunSafe+0x32		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/threading/routines.go:38
#	0x798e8a	github.com/zeromicro/go-zero/core/stat.init.1.func1+0x10a		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/zeromicro/go-zero@v1.10.1/core/stat/usage.go:34

0: 0 [1: 32] @ 0x4099f2 0x40b5b7 0x40b325 0x40c8a9 0x8993e7 0x89a06e 0x898d33 0x898cc3 0x455058 0x446445 0x44632e 0x4835a1
#	0x8993e6	github.com/prometheus/client_golang/prometheus.(*Registry).Register+0x346	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/prometheus/client_golang@v1.23.2/prometheus/registry.go:306
#	0x89a06d	github.com/prometheus/client_golang/prometheus.(*Registry).MustRegister+0x4d	/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/prometheus/client_golang@v1.23.2/prometheus/registry.go:405
#	0x898d32	github.com/prometheus/client_golang/prometheus.MustRegister+0xf2		/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/prometheus/client_golang@v1.23.2/prometheus/registry.go:177
#	0x898cc2	github.com/prometheus/client_golang/prometheus.init.0+0x82			/home/fanqicheng/.gvm/pkgsets/go1.24.4/global/pkg/mod/github.com/prometheus/client_golang@v1.23.2/prometheus/registry.go:62
#	0x455057	runtime.doInit1+0xd7								/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:7353
#	0x446444	runtime.doInit+0x344								/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:7320
#	0x44632d	runtime.main+0x22d								/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:254

1: 2048 [1: 2048] @ 0x44a231 0x44ad55 0x44b439 0x47b76c 0x44d9be 0x44de4f 0x44e2a5 0x48156e
#	0x44a230	runtime.allocm+0x90		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2236
#	0x44ad54	runtime.newm+0x34		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2772
#	0x44b438	runtime.startm+0x158		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2998
#	0x47b76b	runtime.wakep+0xeb		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3145
#	0x44d9bd	runtime.resetspinning+0x3d	/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3885
#	0x44de4e	runtime.schedule+0x10e		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:4038
#	0x44e2a4	runtime.park_m+0x284		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:4144
#	0x48156d	runtime.mcall+0x4d		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/asm_amd64.s:459

3: 6144 [3: 6144] @ 0x44a231 0x44ad55 0x44b439 0x47b76c 0x44d9be 0x44de4f 0x44996d 0x449875 0x4814e5
#	0x44a230	runtime.allocm+0x90		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2236
#	0x44ad54	runtime.newm+0x34		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2772
#	0x44b438	runtime.startm+0x158		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:2998
#	0x47b76b	runtime.wakep+0xeb		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3145
#	0x44d9bd	runtime.resetspinning+0x3d	/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:3885
#	0x44de4e	runtime.schedule+0x10e		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:4038
#	0x44996c	runtime.mstart1+0xcc		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:1862
#	0x449874	runtime.mstart0+0x74		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/proc.go:1808
#	0x4814e4	runtime.mstart+0x4		/home/fanqicheng/.gvm/gos/go1.24.4/src/runtime/asm_amd64.s:395

1: 6528 [1: 6528] @ 0x51a765 0x5215a5 0x49936b 0x51ab25 0x51aaf4 0xa6f53f 0xa6f51d 0xa6f2fd 0x764a49 0x766944 0x78558e 0x762f45 0x4835a1
#	0x51a764	strings.(*Replacer).build+0x164		/home/fanqicheng/.gvm/gos/go1.24.4/src/strings/replace.go:75
#	0x5215a4	strings.(*Replacer).buildOnce+0x24	/home/fanqicheng/.gvm/gos/go1.24.4/src/strings/replace.go:40
#	0x49936a	sync.(*Once).doSlow+0xaa		/home/fanqicheng/.gvm/gos/go1.24.4/src/sync/once.go:78
#	0x51ab24	sync.(*Once).Do+0x44			/home/fanqicheng/.gvm/gos/go1.24.4/src/sync/once.go:69
#	0x51aaf3	strings.(*Replacer).Replace+0x13	/home/fanqicheng/.gvm/gos/go1.24.4/src/strings/replace.go:96
#	0xa6f53e	html.EscapeString+0xbe			/home/fanqicheng/.gvm/gos/go1.24.4/src/html/escape.go:179
#	0xa6f51c	net/http/pprof.indexTmplExecute+0x9c	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/pprof/pprof.go:449
#	0xa6f2fc	net/http/pprof.Index+0x73c		/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/pprof/pprof.go:420
#	0x764a48	net/http.HandlerFunc.ServeHTTP+0x28	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2294
#	0x766943	net/http.(*ServeMux).ServeHTTP+0x1c3	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2822
#	0x78558d	net/http.serverHandler.ServeHTTP+0x8d	/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:3301
#	0x762f44	net/http.(*conn).serve+0x624		/home/fanqicheng/.gvm/gos/go1.24.4/src/net/http/server.go:2102


# runtime.MemStats
# Alloc = 2571712
# TotalAlloc = 32346456
# Sys = 13980936
# Lookups = 0
# Mallocs = 140035
# Frees = 128217
# HeapAlloc = 2571712
# HeapSys = 7700480
# HeapIdle = 3710976
# HeapInuse = 3989504
# HeapReleased = 1720320
# HeapObjects = 11818
# Stack = 688128 / 688128
# MSpan = 113760 / 146880
# MCache = 9664 / 15704
# BuckHashSys = 1449349
# GCSys = 2768848
# OtherSys = 1211547
# NextGC = 4194304
# LastGC = 1779903435095809950
# PauseNs = [731656 123293 215062 288174 90099 111835 173989 61312 229928 524563 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
# PauseEnd = [1779902971096933709 1779903023093746511 1779903075093836922 1779903123094353344 1779903175092891904 1779903227345310578 1779903279844904107 1779903331092838078 1779903383594286723 1779903435095809950 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
# NumGC = 10
# NumForcedGC = 0
# GCCPUFraction = 9.093591992995519e-06
# DebugGC = false
# MaxRSS = 42393600