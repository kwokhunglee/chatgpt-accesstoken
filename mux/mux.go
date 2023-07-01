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

package mux

import (
	"net/http"

	akt "github.com/chatgpt-accesstoken"
	"github.com/gin-gonic/gin"
)

type Server struct {
	openAuthSvc akt.OpenaiAuthService
}

func New(openAuthSvc akt.OpenaiAuthService) *Server {
	return &Server{
		openAuthSvc: openAuthSvc,
	}
}

func (s Server) Handler() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Any("/health", Healthy)

	r.Group("/auth", func(context *gin.Context) {
		r.POST("/", handlerPostAuth)
		r.POST("/puid", handlerPostPUID)
		r.POST("/all", handlerPostAll)
	})

	pg := r.Group("/proxy")
	{
		pg.GET("/", handlerGetProxy)
		pg.POST("/", handlerPostProxy)
		pg.DELETE("/", handlerDeleteProxy)
	}
	return r
}

func Healthy(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func handlerPostAuth(ctx *gin.Context) {

}

func handlerPostPUID(ctx *gin.Context) {

}

func handlerPostAll(ctx *gin.Context) {

}

func handlerGetProxy(ctx *gin.Context) {

}

func handlerPostProxy(ctx *gin.Context) {

}

func handlerDeleteProxy(ctx *gin.Context) {

}
