/*
Copyright 2022 The Workpieces LLC.

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

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/workpieces/log"
	"net/http"
)

func Healthy(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func Auth(ctx *gin.Context) {

}

func main() {
	logg := log.NewNop()

	logg.WithField("port", "")

	g := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	g.Any("/health", Healthy)
	g.POST("/auth", Auth)

	if err := g.Run(); err != nil {
		panic(fmt.Errorf("http server run err: %s", err))
	}
}
