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

package models

type PublicRoomEntry struct {
	CanonicalAlias string           `json:"canonical_alias,omitempty"`
	Name           string           `json:"name,omitempty"`
	JoinedCount    int              `json:"num_joined_members"`
	RoomID         string           `json:"room_id"`
	Topic          string           `json:"topic,omitempty"`
	WorldReadable  bool             `json:"world_readable"`
	GuestsAllowed  bool             `json:"guest_can_join"`
	JoinRule       string           `json:"join_rule,omitempty"`
	RoomType       string           `json:"room_type"`
	AvatarUrl      string           `json:"avatar_url,omitempty"`
	ChildrenState  []*ChildrenState `json:"children_state"`
}

type ChildrenState struct {
	Type           string                 `json:"type"`
	StateKey       string                 `json:"state_key,omitempty"`
	Sender         string                 `json:"sender"`
	OriginServerTs int64                  `json:"origin_server_ts"`
	Content        map[string]interface{} `json:"content"`
}
