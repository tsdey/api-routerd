// SPDX-License-Identifier: Apache-2.0

package main

import (
	"api-routerd/cmd/router"
	"api-routerd/cmd/share"
	log "github.com/sirupsen/logrus"
	"github.com/go-ini/ini"
	"flag"
	"runtime"
)

// Version app version
const Version = "0.1"
const ConfPath = "/etc/api-routerd.conf"

var ipFlag string
var portFlag string
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

func InitConf() {
	cfg, r := ini.Load(ConfPath)
	if r != nil {
		log.Errorf("Fail to read conf file '%s': %v", ConfPath, r)
		return
	}

	ip := cfg.Section("Network").Key("IPAddress").String()
	port := cfg.Section("Network").Key("Port").String()

	log.Infof("Conf file IPAddress=%s, Port=%s", ip, port)

	if ip != "" {
		ipFlag = ip
	}

	if port != "" {
		portFlag = port
	}
}

func main() {
	share.InitLog()
	InitConf()
	flag.Parse()

	log.Infof("api-routerd: v%s (built %s)", Version, runtime.Version())
	log.Infof("Start Server at %s:%s", ipFlag, portFlag)

	router.StartRouter(ipFlag, portFlag)
}
