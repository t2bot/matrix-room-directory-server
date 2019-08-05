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
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/api/common"
	"github.com/t2bot/matrix-room-directory-server/directory"
	"github.com/t2bot/matrix-room-directory-server/key_server"
	"github.com/t2bot/matrix-room-directory-server/util"
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

	rooms, err := directory.Default.GetRooms()
	if err != nil {
		log.Error(err)
		return common.InternalServerError("failed to get room list")
	}

	limitRaw := r.URL.Query().Get("limit")
	sinceRaw := r.URL.Query().Get("since")

	limit := 0
	since := 0
	if limitRaw != "" {
		v, err := strconv.Atoi(limitRaw)
		if err != nil {
			log.Error(err)
			return common.InternalServerError("failed to parse limit")
		}
		limit = v
	}
	if sinceRaw != "" {
		v, err := strconv.Atoi(sinceRaw)
		if err != nil {
			log.Error(err)
			return common.InternalServerError("failed to parse since")
		}
		since = v
	}

	max := len(rooms)
	start := since
	end := util.Min(max, start+limit)
	if end == start {
		end = max
	}

	subsetRooms := rooms[start:end]

	nextToken := ""
	if end != max {
		nextToken = strconv.Itoa(end + 1)
	}

	prevToken := ""
	if start != 0 {
		prevToken = strconv.Itoa(start - 1)
	}

	chunk := make([]*PublicRoomEntry, 0)
	for _, r := range subsetRooms {
		chunk = append(chunk, &PublicRoomEntry{
			RoomID:         r.RoomID,
			WorldReadable:  r.WorldReadable,
			AvatarUrl:      r.AvatarUrl,
			Topic:          r.Topic,
			Name:           r.Name,
			CanonicalAlias: r.CanonicalAlias,
			GuestsAllowed:  r.GuestsCanJoin,
			JoinedCount:    r.JoinedMembers,
			Aliases:        []string{r.CanonicalAlias}, // TODO: Track aliases correctly
		})
	}

	return &PublicRoomsResponse{
		Chunk:           chunk,
		NextBatchToken:  nextToken,
		PrevBatchToken:  prevToken,
		TotalRoomsKnown: len(rooms),
	}
}
