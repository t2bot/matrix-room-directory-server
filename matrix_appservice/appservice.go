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

package matrix_appservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/common"
	"github.com/t2bot/matrix-room-directory-server/util"
)

type MinimalMatrixEvent struct {
	EventID   string                 `json:"event_id"`
	RoomID    string                 `json:"room_id"`
	EventType string                 `json:"type"`
	Sender    string                 `json:"sender"`
	StateKey  string                 `json:"state_key,omitempty"`
	Content   map[string]interface{} `json:"content"`
}

type Appservice struct {
	homeserverUrl string
	asToken       string
	hsToken       string
	userId        string
}

var Default *Appservice

func Setup(hsUrl string, asToken string, hsToken string) error {
	app, err := NewAppservice(hsUrl, asToken, hsToken)
	if err != nil {
		return err
	}
	logrus.Info("Default appservice running as ", app.userId)
	Default = app
	return nil
}

func NewAppservice(hsUrl string, asToken string, hsToken string) (*Appservice, error) {
	app := &Appservice{homeserverUrl: hsUrl, asToken: asToken, hsToken: hsToken}
	userId, err := app.GetUserId()
	if err != nil {
		return nil, err
	}
	app.userId = userId
	return app, nil
}

func (a *Appservice) GetUserId() (string, error) {
	r, err := a.doRequest("GET", "/_matrix/client/r0/account/whoami", nil)
	if err != nil {
		return "", err
	}

	return r["user_id"].(string), nil
}

func (a *Appservice) CheckHomeserverAuth(r *http.Request) error {
	token := util.GetAccessTokenFromRequest(r)
	if token == "" {
		return errors.New("no token found")
	}

	if token != a.hsToken {
		return errors.New("token does not match")
	}

	return nil
}

func (a *Appservice) JoinRoom(roomIdOrAlias string) (string, error) {
	r, err := a.doRequest("POST", "/_matrix/client/r0/join/"+url.QueryEscape(roomIdOrAlias), nil)
	if err != nil {
		return "", err
	}

	return r["room_id"].(string), nil
}

func (a *Appservice) ProcessEvent(ev *MinimalMatrixEvent) error {
	if ev.Content == nil {
		return nil
	}

	if ev.EventType == "m.room.member" && ev.StateKey == a.userId && ev.Sender == common.AdminUser && ev.Content["membership"].(string) == "invite" {
		logrus.Info("Received invite from admin")
		_, err := a.JoinRoom(ev.RoomID)
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

func (a *Appservice) doRequest(method string, path string, body interface{}) (map[string]interface{}, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, a.homeserverUrl+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+a.asToken)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(r.Body)
		if b != nil && len(b) > 0 {
			logrus.Error(string(b))
		}
		return nil, errors.New("request failed")
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
