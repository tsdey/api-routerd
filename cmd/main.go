// SPDX-License-Identifier: Apache-2.0

package main

import (
	"api-routerd/cmd/router"
	"api-routerd/cmd/share"
	log "github.com/sirupsen/logrus"
	"flag"
	"runtime"
)

// Version app version
const Version = "0.1"

var portFlag string
var ipFlag string
var tokenFlag string

func init() {
	const (
		defaultIP  = "0.0.0.0"
		defaultPort  = "8080"
		defaultToken = "token.txt"
	)

	flag.StringVar(&ipFlag, "ip", defaultIP, "The server IP address.")
	flag.StringVar(&portFlag, "port", defaultPort, "The server port.")
	flag.StringVar(&tokenFlag, "token", defaultToken, "The token file for authentication.")
}

func main() {
	share.InitLog()
	flag.Parse()

	log.Infof("api-routerd: v%s (built %s)", Version, runtime.Version())
	log.Infof("Start Server at %s:%s", ipFlag, portFlag)

	router.StartRouter(ipFlag, portFlag)
}
