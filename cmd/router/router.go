// SPDX-License-Identifier: Apache-2.0

package router

import (
	"api-routerd/cmd/hostname"
	"api-routerd/cmd/network"
	"api-routerd/cmd/proc"
	"api-routerd/cmd/systemd"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartRouter(ip string, port string) {
	router := mux.NewRouter()

	// Register services
	hostname.RegisterRouterHostname(router)
	network.RegisterRouterNetwork(router)
	proc.RegisterRouterProc(router)
	systemd.RegisterRouterSystemd(router)

	log.Fatal(http.ListenAndServe(ip + ":" + port, router))
}
