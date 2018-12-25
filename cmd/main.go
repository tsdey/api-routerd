// SPDX-License-Identifier: Apache-2.0

package main

import (
	"api-routerd/cmd/router"
	"api-routerd/cmd/share"
	log "github.com/sirupsen/logrus"
	"github.com/go-ini/ini"
	"flag"
	"runtime"
	"os"
	"path"
)

// Version app version
const Version = "0.1"
const ConfPath = "/etc/api-routerd"
const ConfFile = "api-routerd.conf"
const TlsCert = "tls/server.crt"
const TlsKey = "tls/server.key"

var ipFlag string
var portFlag string

func init() {
	const (
		defaultIP  = "0.0.0.0"
		defaultPort  = "8080"
	)

	flag.StringVar(&ipFlag, "ip", defaultIP, "The server IP address.")
	flag.StringVar(&portFlag, "port", defaultPort, "The server port.")
}

func InitConf() {
	confFile := path.Join(ConfPath, ConfFile)
	cfg, r := ini.Load(confFile)
	if r != nil {
		log.Errorf("Fail to read conf file '%s': %v", ConfPath, r)
		return
	}

	ip := cfg.Section("Network").Key("IPAddress").String()
	port := cfg.Section("Network").Key("Port").String()

	log.Infof("Conf file IPAddress=%s, Port=%s", ip, port)

	if ip != "" && port != "" {
		ipFlag = ip
		portFlag = port
	}
}

func main() {
	share.InitLog()
	InitConf()
	flag.Parse()

	log.Infof("api-routerd: v%s (built %s)", Version, runtime.Version())
	log.Infof("Start Server at %s:%s", ipFlag, portFlag)

	r := router.StartRouter(ipFlag, portFlag, path.Join(ConfPath, TlsCert), path.Join(ConfPath, TlsKey))
	if r != nil {
		log.Fatal("Failed to init router: %s", r)
		os.Exit(1)
	}
}
