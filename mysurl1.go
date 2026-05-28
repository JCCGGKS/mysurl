// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"

	"mysurl1/internal/config"
	"mysurl1/internal/handler"
	"mysurl1/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/mysurl1-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if c.Stat.DisableLog {
		stat.DisableLog()
	}

	if !c.Stat.DisableSampler {
		// This branch is the config-controlled entry point. Fully preventing the
		// sampler from starting still requires moving the start logic out of the
		// go-zero package init path.
		stat.Stat()
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
