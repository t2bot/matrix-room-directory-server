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

package appservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/api/common"
	"github.com/t2bot/matrix-room-directory-server/matrix_appservice"
)

type Transaction struct {
	Events []*matrix_appservice.MinimalMatrixEvent `json:"events"`
}

func ReceiveTransaction(r *http.Request, log *logrus.Entry) interface{} {
	err := matrix_appservice.Default.CheckHomeserverAuth(r)
	if err != nil {
		log.Error(err)
		return common.UnauthorizedError()
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return common.InternalServerError("Failed to read body")
	}

	var txn Transaction
	err = json.Unmarshal(b, &txn)
	if err != nil {
		log.Error(err)
		return common.BadRequest("Not JSON")
	}

	for _, ev := range txn.Events {
		err = matrix_appservice.Default.ProcessEvent(ev)
		if err != nil {
			log.Error(err)
			return common.InternalServerError("failed to process event " + ev.EventID)
		}
	}

	return &common.EmptyResponse{}
}
