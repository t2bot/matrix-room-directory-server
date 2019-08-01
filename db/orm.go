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

	"github.com/t2bot/matrix-room-directory-server/db/models"
)

func GetRooms() ([]*models.DirectoryRoom, error) {
	r, err := statements[selectAllRooms].Query()
	if err == sql.ErrNoRows {
		return make([]*models.DirectoryRoom, 0), nil
	}
	if err != nil {
		return nil, err
	}

	var results []*models.DirectoryRoom
	for r.Next() {
		v := &models.DirectoryRoom{}
		err = r.Scan(&v.RoomID, &v.CanonicalAlias, &v.Name, &v.Topic, &v.AvatarUrl, &v.JoinedMembers, &v.WorldReadable, &v.GuestsCanJoin)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func UpsertRoom(roomId string, canonicalAlias string, name string, topic string, avatarUrl string, joinedMembers int, worldReadable bool, guestsCanJoin bool) error {
	_, err := statements[upsertRoom].Exec(roomId, canonicalAlias, name, topic, avatarUrl, joinedMembers, worldReadable, guestsCanJoin)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRoom(roomId string) error {
	_, err := statements[deleteRoom].Exec(roomId)
	if err != nil {
		return err
	}
	return nil
}
