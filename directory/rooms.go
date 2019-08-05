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

package directory

import (
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/db"
	"github.com/t2bot/matrix-room-directory-server/db/models"
	"github.com/t2bot/matrix-room-directory-server/matrix_appservice"
)

var Default *Directory

type Directory struct {
	appservice  *matrix_appservice.Appservice
	cachedRooms []*models.DirectoryRoom
}

func New(appservice *matrix_appservice.Appservice) *Directory {
	return &Directory{appservice: appservice}
}

func (d *Directory) AddRoom(roomIdOrAlias string) error {
	roomId, err := d.appservice.JoinRoom(roomIdOrAlias)
	if err != nil {
		return err
	}

	// Update also adds the room
	return d.UpdateRoom(roomId)
}

func (d *Directory) UpdateRoom(roomId string) error {
	logrus.Info("Updating room: " + roomId)

	state, err := d.appservice.GetRoomState(roomId)
	if err != nil {
		return err
	}

	canonicalAlias := ""
	name := ""
	topic := ""
	avatarUrl := ""
	joinedCount := 0
	worldReadable := false
	guestsCanJoin := false
	isPublic := false

	for _, s := range state {
		if s.EventType == "m.room.name" && s.StateKey == "" {
			v, ok := s.Content["name"]
			if !ok {
				name = ""
			} else {
				name = v.(string)
			}
		} else if s.EventType == "m.room.avatar" && s.StateKey == "" {
			v, ok := s.Content["url"]
			if !ok {
				avatarUrl = ""
			} else {
				avatarUrl = v.(string)
			}
		} else if s.EventType == "m.room.topic" && s.StateKey == "" {
			v, ok := s.Content["topic"]
			if !ok {
				topic = ""
			} else {
				topic = v.(string)
			}
		} else if s.EventType == "m.room.history_visibility" && s.StateKey == "" {
			v, ok := s.Content["history_visibility"]
			if !ok {
				worldReadable = false
			} else {
				worldReadable = v.(string) == "world_readable"
			}
		} else if s.EventType == "m.room.canonical_alias" && s.StateKey == "" {
			v, ok := s.Content["alias"]
			if !ok {
				canonicalAlias = ""
			} else {
				canonicalAlias = v.(string)
			}
		} else if s.EventType == "m.room.guest_access" && s.StateKey == "" {
			v, ok := s.Content["guest_access"]
			if !ok {
				guestsCanJoin = false
			} else {
				guestsCanJoin = v.(string) == "can_join"
			}
		} else if s.EventType == "m.room.join_rules" && s.StateKey == "" {
			v, ok := s.Content["join_rule"]
			if !ok {
				isPublic = false
			} else {
				isPublic = v.(string) == "public"
			}
		} else if s.EventType == "m.room.member" && s.StateKey != "" {
			m, ok := s.Content["membership"]
			if ok && m.(string) == "join" {
				joinedCount += 1
			}
		}
	}

	viable := isPublic && canonicalAlias != "" && name != ""
	if !viable {
		logrus.Info("Removing room from directory: now private")
		_ = d.appservice.SendNotice(roomId, "This room has been removed from the directory because it is now private")
		err = db.DeleteRoom(roomId)
		d.cachedRooms = nil
		if err != nil {
			return err
		}
		return nil
	}

	err = db.UpsertRoom(roomId, canonicalAlias, name, topic, avatarUrl, joinedCount, worldReadable, guestsCanJoin)
	if err != nil {
		return err
	}

	d.cachedRooms = nil
	return nil
}

func (d *Directory) GetRooms() ([]*models.DirectoryRoom, error) {
	if d.cachedRooms != nil {
		return d.cachedRooms, nil
	}

	rooms, err := db.GetRooms()
	if err != nil {
		return nil, err
	}

	d.cachedRooms = rooms
	return rooms, nil
}
