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

package federation

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/api/common"
	"github.com/t2bot/matrix-room-directory-server/key_server"
)

type PublicRoomEntry struct {
	Aliases        []string `json:"aliases"`
	CanonicalAlias string   `json:"canonical_alias,omitempty"`
	Name           string   `json:"name,omitempty"`
	JoinedCount    int      `json:"num_joined_members"`
	RoomID         string   `json:"room_id"`
	Topic          string   `json:"topic,omitempty"`
	WorldReadable  bool     `json:"world_readable"`
	GuestsAllowed  bool     `json:"guest_can_join"`
	AvatarUrl      string   `json:"avatar_url,omitempty"`
}

type PublicRoomsResponse struct {
	Chunk           []*PublicRoomEntry `json:"chunk"`
	NextBatchToken  string             `json:"next_batch,omitempty"`
	PrevBatchToken  string             `json:"prev_batch,omitempty"`
	TotalRoomsKnown int                `json:"total_room_count_estimate"`
}

func GetPublicRooms(r *http.Request, log *logrus.Entry) interface{} {
	auth := r.Header.Get("Authorization")
	urlWithQuery := r.URL.Path + "?" + r.URL.RawQuery
	destination := r.Host
	method := r.Method

	err := key_server.Default.CheckAuth(auth, method, urlWithQuery, destination)
	if err != nil {
		log.Error(err)
		return common.InternalServerError("failed to authenticate request or some other error")
	}

	return &PublicRoomsResponse{
		Chunk: []*PublicRoomEntry{
			{
				Aliases:        []string{"#matrix:matrix.org"},
				CanonicalAlias: "#matrix:matrix.org",
				Name:           "Matrix HQ [TEST]",
				JoinedCount:    1112,
				RoomID:         "!QtykxKocfZaZOUrTwp:matrix.org",
				Topic:          "Testing the directory server",
				WorldReadable:  true,
				GuestsAllowed:  true,
				AvatarUrl:      "mxc://matrix.org/DRevoaEiuzbkOznknySKuMmE",
			},
		},
		TotalRoomsKnown: 1,
	}
}
