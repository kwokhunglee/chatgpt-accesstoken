/*
Copyright 2022 The deepauto-io LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package launcher

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/chatgpt-accesstoken/build"
	"github.com/chatgpt-accesstoken/core"
	mux2 "github.com/chatgpt-accesstoken/mux"
	"github.com/chatgpt-accesstoken/store/redisdb"

	"github.com/workpieces/log"
)

type labeledCloser struct {
	label  string
	closer func(context.Context) error
}

// Launcher represents the main program execution.
type Launcher struct {
	wg       sync.WaitGroup
	cancel   func()
	doneChan <-chan struct{}
	closers  []labeledCloser

	httpPort int
	logger   log.Logger
}

// NewLauncher returns a new instance of Launcher with a no-op logger.
func NewLauncher() *Launcher {
	return &Launcher{
		logger: log.NewNop(),
	}
}

// Shutdown shuts down the HTTP server and waits for all services to clean up.
func (m *Launcher) Shutdown(ctx context.Context) error {
	var errs []string

	// Shut down subsystems in the reverse order of their registration.
	for i := len(m.closers); i > 0; i-- {
		lc := m.closers[i-1]
		m.logger.
			WithField("subsystem", lc.label).
			Info("stopping subsystem")
		if err := lc.closer(ctx); err != nil {
			m.logger.
				WithField("subsystem", lc.label).
				Error("Failed to stop subsystem: ", err)
			errs = append(errs, err.Error())
		}
	}

	m.wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("failed to shut down server: [%s]", strings.Join(errs, ","))
	}
	return nil
}

func (m *Launcher) Done() <-chan struct{} {
	return m.doneChan
}

func (m *Launcher) run(ctx context.Context, opts Config) error {
	ctx, m.cancel = context.WithCancel(ctx)
	m.doneChan = ctx.Done()

	info := build.GetBuildInfo()

	m.logger.
		WithField("version", info.Version).
		WithField("commit", info.Commit).
		WithField("build_date", info.Date).
		Info("welcome to akt")

	db := redisdb.New(opts.RedisDB)
	m.closers = append(m.closers, labeledCloser{
		label: "Redis Server",
		closer: func(ctx context.Context) error {
			return db.Close()
		},
	})

	openaiAuthSvc := core.NewOpenaiAuthLogger(m.logger,
		core.NewOpenaiAuthCache(db,
			core.New()))

	srv := &http.Server{
		Addr:    opts.HttpBindAddress,
		Handler: mux2.New(openaiAuthSvc, core.NewProxyService(db)).Handler(),
	}

	m.closers = append(m.closers, labeledCloser{
		label:  "HTTP Server",
		closer: srv.Shutdown,
	})

	m.wg.Add(1)
	go func(log log.Logger) {
		defer m.wg.Done()
		log.WithField("port", opts.HttpBindAddress).Info("listening")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("server", "http").Error(err)
			m.cancel()
		}
	}(m.logger)

	return nil
}
