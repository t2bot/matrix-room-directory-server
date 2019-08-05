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

package matrix_appservice_processor

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/common"
	"github.com/t2bot/matrix-room-directory-server/db"
	"github.com/t2bot/matrix-room-directory-server/directory"
	"github.com/t2bot/matrix-room-directory-server/matrix_appservice"
)

var Default *Processor

var directoryEventTypes = map[string]bool{
	"m.room.member":             true,
	"m.room.join_rules":         true,
	"m.room.guest_access":       true,
	"m.room.history_visibility": true,
	"m.room.name":               true,
	"m.room.topic":              true,
	"m.room.avatar":             true,
	"m.room.canonical_alias":    true,
}

type Processor struct {
	appservice   *matrix_appservice.Appservice
	directory    *directory.Directory
	knownRoomIds []string
}

func New(appservice *matrix_appservice.Appservice, dir *directory.Directory) *Processor {
	return &Processor{appservice: appservice, directory: dir}
}

func (a *Processor) ProcessEvent(ev *matrix_appservice.MinimalMatrixEvent) error {
	if ev.Content == nil {
		return nil
	}

	if ev.EventType == "m.room.member" && ev.StateKey == a.appservice.UserID && ev.Sender == common.AdminUser && ev.Content["membership"].(string) == "invite" {
		logrus.Info("Received invite from admin")
		_, err := a.appservice.JoinRoom(ev.RoomID)
		if err != nil {
			logrus.Error(err)
		}
	} else if ev.EventType == "m.room.message" && ev.Sender == common.AdminUser && ev.Content["msgtype"].(string) == "m.text" {
		command := strings.TrimSpace(ev.Content["body"].(string))
		if strings.HasPrefix(command, "!directory add ") {
			a.knownRoomIds = make([]string, 0) // reset cache
			err := a.directory.AddRoom(strings.TrimPrefix(command, "!directory add "))
			if err != nil {
				logrus.Error(err)
				_ = a.appservice.SendReaction(ev.RoomID, ev.EventID, "❌")
			} else {
				err = a.appservice.SendReaction(ev.RoomID, ev.EventID, "✔")
				if err != nil {
					logrus.Error(err)
					_ = a.appservice.SendReaction(ev.RoomID, ev.EventID, "❌")
				}
			}
		} else {
			_ = a.appservice.SendReaction(ev.RoomID, ev.EventID, "❓")
		}
	} else if v, ok := directoryEventTypes[ev.EventType]; ok && v {
		if len(a.knownRoomIds) == 0 {
			rooms, err := db.GetRooms()
			if err != nil {
				logrus.Error(err)
			} else {
				a.knownRoomIds = make([]string, 0)
				for _, r := range rooms {
					a.knownRoomIds = append(a.knownRoomIds, r.RoomID)
				}
			}
		}

		hasRoom := false
		for _, r := range a.knownRoomIds {
			if r == ev.RoomID {
				hasRoom = true
				break
			}
		}

		if hasRoom {
			err := a.directory.UpdateRoom(ev.RoomID)
			if err != nil {
				logrus.Error(err)
			} else {
				_ = a.appservice.SendReaction(ev.RoomID, ev.EventID, "✔")
			}
		}
	}

	return nil
}
