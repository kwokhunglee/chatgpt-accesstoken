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

package core

import (
	"context"
	"math/rand"

	"github.com/acheong08/OpenAIAuth/auth"
	akt "github.com/chatgpt-accesstoken"
)

type openaiAuthCache struct {
	proxySvc akt.ProxyService
	svc      akt.OpenaiAuthService
}

func (o openaiAuthCache) All(ctx context.Context, req *akt.OpenaiAuthRequest) (*auth.AuthResult, error) {
	if req.Proxy == "" {
		list, err := o.proxySvc.List(ctx)
		if err != nil {
			return nil, err
		}

		idx := rand.Intn(len(list))
		req.Proxy = list[idx]
	}
	return o.svc.All(ctx, req)
}

func (o openaiAuthCache) AccessToken(ctx context.Context, req *akt.OpenaiAuthRequest) (*auth.AuthResult, error) {

	if req.Proxy == "" {
		list, err := o.proxySvc.List(ctx)
		if err != nil {
			return nil, err
		}

		idx := rand.Intn(len(list))
		req.Proxy = list[idx]
	}
	return o.svc.AccessToken(ctx, req)
}

func (o openaiAuthCache) PUID(ctx context.Context, req *akt.OpenaiAuthRequest) (*auth.AuthResult, error) {
	if req.Proxy == "" {
		list, err := o.proxySvc.List(ctx)
		if err != nil {
			return nil, err
		}

		idx := rand.Intn(len(list))
		req.Proxy = list[idx]
	}
	return o.svc.PUID(ctx, req)
}

func NewOpenaiAuthCache(proxySvc akt.ProxyService, svc akt.OpenaiAuthService) akt.OpenaiAuthService {
	return &openaiAuthCache{
		proxySvc: proxySvc,
		svc:      svc,
	}
}
