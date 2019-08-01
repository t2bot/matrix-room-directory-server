/*
 * Copyright 2019 Travis Ralston <travis@t2bot.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package key_server

import (
	"errors"
	"net/http"
)

type KeyServer struct {
	url string
}

var Default *KeyServer

func Setup(url string) {
	Default = NewKeyServer(url)
}

func NewKeyServer(url string) *KeyServer {
	return &KeyServer{url}
}

func (k *KeyServer) CheckAuth(authHeader string, urlMethod string, urlWithQuery string, destinationHost string) error {
	contactPath := k.url + "/_matrix/key/unstable/check_auth"

	req, err := http.NewRequest("POST", contactPath, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("X-Keys-Method", urlMethod)
	req.Header.Set("X-Keys-URI", urlWithQuery)
	req.Header.Set("X-Keys-Destination", destinationHost)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New("auth failed")
	}

	return nil
}
