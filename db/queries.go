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

package db

import (
	"database/sql"
)

var statements = map[string]*sql.Stmt{}

const selectAllRooms = "selectAllRooms"
const upsertRoom = "upsertRoom"
const deleteRoom = "deleteRoom"

var queries = map[string]string{
	selectAllRooms: "SELECT room_id, canonical_alias, name, topic, avatar_url, joined_count, world_readable, guests_can_join FROM listed_rooms;",
	upsertRoom:     "INSERT INTO listed_rooms (room_id, canonical_alias, name, topic, avatar_url, joined_count, world_readable, guests_can_join) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (room_id) DO UPDATE SET canonical_alias = $2, name = $3, topic = $4, avatar_url = $5, joined_count = $6, world_readable = $7, guests_can_join = $8;",
	deleteRoom:     "DELETE FROM listed_rooms WHERE room_id = $1;",
}
