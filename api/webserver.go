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

package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-room-directory-server/api/federation"
	"github.com/t2bot/matrix-room-directory-server/api/health"
)

type route struct {
	method  string
	handler handler
}

func Run(listenHost string, listenPort int) {
	rtr := mux.NewRouter()

	healthzHandler := handler{health.Healthz, "healthz"}
	fedPublicRoomsHandler := handler{federation.GetPublicRooms, "federation_public_rooms"}

	routes := make(map[string]route)
	routes["/_matrix/federation/v1/publicRooms"] = route{"GET", fedPublicRoomsHandler}

	for routePath, route := range routes {
		logrus.Info("Registering route: " + route.method + " " + routePath)
		rtr.Handle(routePath, route.handler).Methods(route.method)

		// This is a hack to a ensure that trailing slashes also match the routes correctly
		rtr.Handle(routePath+"/", route.handler).Methods(route.method)
	}

	rtr.Handle("/healthz", healthzHandler).Methods("OPTIONS", "GET")
	rtr.NotFoundHandler = handler{NotFoundHandler, "not_found"}
	rtr.MethodNotAllowedHandler = handler{MethodNotAllowedHandler, "method_not_allowed"}

	address := fmt.Sprintf("%s:%d", listenHost, listenPort)
	httpMux := http.NewServeMux()
	httpMux.Handle("/", rtr)

	logrus.WithField("address", address).Info("Started up. Listening at http://" + address)
	logrus.Fatal(http.ListenAndServe(address, httpMux))
}
