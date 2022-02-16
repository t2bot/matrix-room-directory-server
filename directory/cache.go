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

package directory

import (
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/common"
	"github.com/t2bot/matrix-room-directory-server/matrix"
	"github.com/t2bot/matrix-room-directory-server/models"
	"sort"
	"time"
)

var Cached []*models.PublicRoomEntry

var stopChan = make(chan bool)

func BeginCaching() {
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		defer close(stopChan)
		for {
			select {
			case <-stopChan:
				ticker.Stop()
				return
			case <-ticker.C:
				err := DoUpdate()
				if err != nil {
					logrus.Error("Error updating cache:", err)
				}
			}
		}
	}()
}

func Stop() {
	stopChan <- true
}

func DoUpdate() error {
	logrus.Info("Updating cache...")

	r, err := matrix.GetHierarchy(common.SpaceId)
	if err != nil {
		return err
	}

	// Order the rooms by size
	sort.Slice(r, func(i int, j int) bool {
		return r[i].JoinedCount < r[j].JoinedCount
	})

	Cached = r
	return nil
}
