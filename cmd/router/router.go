// SPDX-License-Identifier: Apache-2.0

package router

import (
	"api-routerd/cmd/hostname"
	"api-routerd/cmd/network"
	"api-routerd/cmd/proc"
	"api-routerd/cmd/systemd"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartRouter(ip string, port string) error {
	router := mux.NewRouter()

	// Register services
	hostname.RegisterRouterHostname(router)
	network.RegisterRouterNetwork(router)
	proc.RegisterRouterProc(router)
	systemd.RegisterRouterSystemd(router)

	// Authenticate users
	amw, r := InitAuthMiddleware()
	if r != nil {
		log.Fatal("Faild to init auth DB existing")
		return errors.New("Failed to init Auth DB")
	}

	router.Use(amw.AuthMiddleware)

	log.Fatal(http.ListenAndServe(ip + ":" + port, router))

	return nil
}
