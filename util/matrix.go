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

package util

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetAccessTokenFromRequest(request *http.Request) string {
	token := request.Header.Get("Authorization")

	if token != "" {
		if !strings.HasPrefix(token, "Bearer") {
			logrus.Warn("Invalid Authorization header observed: expected a Bearer token, got something else")
			return ""
		}
		if len(token) > 7 {
			// "Bearer <token>"
			return token[7:]
		}
	}

	return request.URL.Query().Get("access_token")
}
