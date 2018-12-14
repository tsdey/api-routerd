// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"restgateway/api-router/hostname"
	"restgateway/api-router/network"
	"restgateway/api-router/proc"
	"restgateway/api-router/systemd"
)

func StartRouter() {
	router := mux.NewRouter()

	// Register services
	hostname.RegisterRouterHostname(router)
	network.RegisterRouterNetwork(router)
	proc.RegisterRouterProc(router)
	systemd.RegisterRouterSystemd(router)

	http.ListenAndServe(":8080", router)
}
