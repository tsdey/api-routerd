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

// App Version
const (
	Version  = "0.1"
	ConfPath = "/etc/api-routerd"
	ConfFile = "api-routerd.conf"
	TlsCert  = "tls/server.crt"
	TlsKey   = "tls/server.key"
)

var ipFlag string
var portFlag string

func init() {
	const (
		defaultIP    = ""
		defaultPort  = "8080"
	)

	flag.StringVar(&ipFlag, "ip", defaultIP, "The server IP address.")
	flag.StringVar(&portFlag, "port", defaultPort, "The server port.")
}

func InitConf() {
	confFile := path.Join(ConfPath, ConfFile)
	cfg, err := ini.Load(confFile)
	if err != nil {
		log.Errorf("Fail to read conf file '%s': %v", ConfPath, err)
		return
	}

	ip := cfg.Section("Network").Key("IPAddress").String()
	_, err = share.ParseIP(ip)
	if err != nil {
		log.Errorf("Failed to parse Conf file IPAddress=%s", ip)
		return
	}

	port := cfg.Section("Network").Key("Port").String()
	_, err = share.ParsePort(port)
	if err != nil {
		log.Errorf("Failed to parse Conf file Port=%d", port)
		return
	}

	log.Infof("Conf file IPAddress=%s, Port=%s", ip, port)

	ipFlag = ip
	portFlag = port
}

func main() {
	share.InitLog()
	InitConf()
	flag.Parse()

	log.Infof("api-routerd: v%s (built %s)", Version, runtime.Version())
	log.Infof("Start Server at %s:%s", ipFlag, portFlag)

	err := router.StartRouter(ipFlag, portFlag, path.Join(ConfPath, TlsCert), path.Join(ConfPath, TlsKey))
	if err != nil {
		log.Fatal("Failed to init router: %s", err)
		os.Exit(1)
	}
}
