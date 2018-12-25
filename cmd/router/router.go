// SPDX-License-Identifier: Apache-2.0

package router

import (
	"api-routerd/cmd/hostname"
	"api-routerd/cmd/network"
	"api-routerd/cmd/proc"
	"api-routerd/cmd/share"
	"api-routerd/cmd/systemd"
	"crypto/tls"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartRouter(ip string, port string, tlsCertPath string, tlsKeyPath string) error {
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

	if share.PathExists(tlsCertPath) && share.PathExists(tlsKeyPath) {
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: false,
		}
		srv := &http.Server{
			Addr:         ip + ":" + port,
			Handler:      router,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		log.Info("Starting api-routerd in TLS mode")

		log.Fatal(srv.ListenAndServeTLS(tlsCertPath, tlsKeyPath))

	} else {
		log.Info("Starting api-routerd in plain text mode")

		log.Fatal(http.ListenAndServe(ip + ":" + port, router))
	}

	return nil
}
