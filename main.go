/*
 * Copyright 2019 - 2022 Travis Ralston <travis@t2bot.io>
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

package main

import (
	"github.com/namsral/flag"
	"github.com/t2bot/matrix-room-directory-server/api"
	"github.com/t2bot/matrix-room-directory-server/common"
	"github.com/t2bot/matrix-room-directory-server/directory"
	"github.com/t2bot/matrix-room-directory-server/key_server"
	"github.com/t2bot/matrix-room-directory-server/matrix"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/logging"
)

func main() {
	logging.Setup()
	logrus.Info("Starting up...")

	accessToken := flag.String("accesstoken", "", "Access token to make homeserver API calls with")
	hsUrl := flag.String("hsurl", "https://t2bot.io", "Homeserver to run against")
	keyServerUrl := flag.String("keyserver", "https://keys.t2host.io", "Key server to perform auth against")
	spaceId := flag.String("space", "#directory:t2bot.io", "The Space to use as a room directory")
	listenHost := flag.String("address", "0.0.0.0", "Address to listen for requests on")
	listenPort := flag.Int("port", 8080, "Port to listen for requests on")
	flag.Parse()

	logrus.Info("Setting common variables...")
	common.AccessToken = *accessToken
	common.HomeserverUrl = *hsUrl
	common.SpaceId = *spaceId

	logrus.Info("Homeserver URL: ", common.HomeserverUrl)
	logrus.Info("Space ID: ", common.SpaceId)

	logrus.Info("Resolving Space ID to Room ID...")
	rid, err := matrix.ResolveRoom(common.SpaceId)
	if err != nil {
		panic(err)
	}
	common.SpaceId = rid
	logrus.Info("Space ID (revised): ", common.SpaceId)

	logrus.Info("Seeing cache...")
	err = directory.DoUpdate()
	if err != nil {
		panic(err)
	}

	logrus.Info("Setting up key server...")
	key_server.Setup(*keyServerUrl)

	logrus.Info("Starting app...")
	directory.BeginCaching()
	api.Run(*listenHost, *listenPort)

	logrus.Info("Stopping...")
	directory.Stop()
}
