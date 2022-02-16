/*
 * Copyright 2022 Travis Ralston <travis@t2bot.io>
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

package matrix

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/t2bot/matrix-room-directory-server/common"
	"github.com/t2bot/matrix-room-directory-server/models"
	"io/ioutil"
	"net/http"
	"net/url"
)

type directoryLookupResponse struct {
	RoomId string `json:"room_id"`
}

type spaceHierarchyResponse struct {
	Chunk     []*models.PublicRoomEntry `json:"rooms"`
	NextBatch string                    `json:"next_batch"`
}

func ResolveRoom(roomAlias string) (string, error) {
	if roomAlias[0] == '!' {
		return roomAlias, nil
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/_matrix/client/r0/directory/room/%s", common.HomeserverUrl, url.QueryEscape(roomAlias)), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+common.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("unexpected status code %d", res.StatusCode))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	j := directoryLookupResponse{}
	err = json.Unmarshal(b, &j)
	if err != nil {
		return "", err
	}

	return j.RoomId, nil
}

func GetHierarchy(roomId string) ([]*models.PublicRoomEntry, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/_matrix/client/v1/rooms/%s/hierarchy?limit=1000&max_depth=10", common.HomeserverUrl, url.QueryEscape(roomId)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+common.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("unexpected status code %d", res.StatusCode))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	j := spaceHierarchyResponse{}
	err = json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	return j.Chunk, nil
}
