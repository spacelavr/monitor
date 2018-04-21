package monitor

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spacelavr/monitor/pkg/log"
	"github.com/spacelavr/monitor/pkg/monitor/cri"
	"github.com/spacelavr/monitor/pkg/monitor/env"
	"github.com/spacelavr/monitor/pkg/monitor/metrics"
	"github.com/spacelavr/monitor/pkg/monitor/router"
	"github.com/spf13/viper"
)

func Daemon() {
	log.Debug("start monitor daemon")

	var (
		sig = make(chan os.Signal)
	)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	m := metrics.New()

	go func() {
		if err := m.Collect(); err != nil {
			log.Fatal(err)
		}
	}()

	env.SetCri(cri.New())
	env.SetMetrics(m)

	go func() {
		srv := &http.Server{
			Handler: router.Router(),
			Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-sig
	log.Debug("handle SIGINT and SIGTERM")
}