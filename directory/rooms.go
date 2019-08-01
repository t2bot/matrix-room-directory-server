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
	"github.com/t2bot/matrix-room-directory-server/matrix_appservice"
)

func AddRoom(roomIdOrAlias string) error {
	as := matrix_appservice.Default
	roomId, err := as.JoinRoom(roomIdOrAlias)
	if err != nil {
		return err
	}

	// Update also adds the room
	return UpdateRoom(roomId)
}

func UpdateRoom(roomId string) error {
	// TODO: get all state events, add to directory table
	return nil
}
