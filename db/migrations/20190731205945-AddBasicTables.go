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

package migrations

import (
	"database/sql"
)

func Up20190731205945AddBasicTables(db *sql.DB) error {
	var err error

	_, err = db.Exec("CREATE TABLE listed_rooms (room_id VARCHAR(255) NOT NULL PRIMARY KEY, canonical_alias VARCHAR(255) NULL, name VARCHAR(255) NULL, joined_count INT NOT NULL, topic VARCHAR(1024) NULL, world_readable BOOLEAN NOT NULL DEFAULT FALSE, guests_can_join BOOLEAN NOT NULL DEFAULT FALSE, avatar_url VARCHAR(255) NULL);")
	if err != nil {
		return err
	}

	return nil
}
