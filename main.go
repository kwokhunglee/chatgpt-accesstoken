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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/chatgpt"
	"github.com/workpieces/log"
)

type handler struct {
	logger log.Logger
}

func NewHandler(logger log.Logger) *handler {
	return &handler{logger: logger}
}

func (h *handler) Healthy(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

type accessTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type accessTokenResponse struct {
	Expires      time.Time `json:"expires"`
	AccessToken  string    `json:"accessToken"`
	AuthProvider string    `json:"authProvider"`
}

var tokensMap = make(map[string]*accessTokenResponse)

func (h *handler) Auth(ctx *gin.Context) {
	in := new(accessTokenRequest)
	if err := ctx.BindJSON(in); err != nil {
		h.logger.WithField("err", err).Error("api: cannot bind json")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	llog := h.logger.
		WithField("email", in.Email).
		WithField("password", in.Password)

	if ak := tokensMap[in.Email]; ak != nil {
		expire := ak.Expires
		if time.Now().Before(expire) {
			llog.Info("api: token not expire")
			ctx.JSON(http.StatusOK, gin.H{"default": ak.AccessToken}) // support 潘多拉
			return
		} else {
			llog.Info("api: token has expire")
			delete(tokensMap, in.Email)
			goto LABEL
		}
	}

LABEL:
	resp, err := chatgpt.GetAccessToken(api.LoginInfo{
		Username: in.Email,
		Password: in.Password,
	})
	if err != nil {
		llog.WithField("err", err).Error("api: cannot get access_token")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	var ak accessTokenResponse
	if err := json.Unmarshal(resp, &ak); err != nil {
		llog.WithField("err", err).Info("api: cannot json unmarshal")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	llog.Info("token does not exist")
	tokensMap[in.Email] = &ak
	ctx.JSON(http.StatusOK, gin.H{"default": ak.AccessToken}) // support 潘多拉
}

func main() {
	logg := log.NewNop()
	plog := logg.
		WithField("Proxy", os.Getenv("GO_CHATGPT_API_PROXY")).
		WithField("port", "8080")

	h := NewHandler(plog)

	g := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	g.Any("/health", h.Healthy)
	g.POST("/auth", h.Auth)

	if err := g.Run(); err != nil {
		panic(fmt.Errorf("http server run err: %s", err))
	}
}
