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
	"time"

	"github.com/sirupsen/logrus"
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
	UserID        string
}

var Default *Appservice

func Setup(hsUrl string, asToken string, hsToken string) error {
	app, err := NewAppservice(hsUrl, asToken, hsToken)
	if err != nil {
		return err
	}
	logrus.Info("Default appservice running as ", app.UserID)
	Default = app
	return nil
}

func NewAppservice(hsUrl string, asToken string, hsToken string) (*Appservice, error) {
	app := &Appservice{homeserverUrl: hsUrl, asToken: asToken, hsToken: hsToken}
	userId, err := app.GetUserId()
	if err != nil {
		return nil, err
	}
	app.UserID = userId
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

func (a *Appservice) GetRoomState(roomId string) ([]*MinimalMatrixEvent, error) {
	r, err := a.doArrayRequest("GET", "/_matrix/client/r0/rooms/"+url.QueryEscape(roomId)+"/state", nil)
	if err != nil {
		return nil, err
	}

	enc, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var result []*MinimalMatrixEvent
	err = json.Unmarshal(enc, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Appservice) SendReaction(roomId string, eventId string, reaction string) error {
	body := make(map[string]interface{})
	relationship := make(map[string]interface{})
	relationship["rel_type"] = "m.annotation"
	relationship["event_id"] = eventId
	relationship["key"] = reaction
	body["m.relates_to"] = relationship

	_, err := a.doRequest("PUT", "/_matrix/client/r0/rooms/"+url.QueryEscape(roomId)+"/send/m.reaction/"+time.Now().String(), body)
	if err != nil {
		return err
	}
	return nil
}

func (a *Appservice) SendNotice(roomId string, notice string) error {
	body := make(map[string]interface{})
	body["msgtype"] = "m.notice"
	body["body"] = notice

	_, err := a.doRequest("PUT", "/_matrix/client/r0/rooms/"+url.QueryEscape(roomId)+"/send/m.room.message/"+time.Now().String(), body)
	if err != nil {
		return err
	}
	return nil
}

func (a *Appservice) doRequest(method string, path string, body interface{}) (map[string]interface{}, error) {
	b, err := a.doRawRequest(method, path, body)
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

func (a *Appservice) doArrayRequest(method string, path string, body interface{}) ([]interface{}, error) {
	b, err := a.doRawRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	var m []interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (a *Appservice) doRawRequest(method string, path string, body interface{}) ([]byte, error) {
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

	return b, nil
}
