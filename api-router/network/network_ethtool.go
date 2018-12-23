package network

import (
	"encoding/json"
	"github.com/safchain/ethtool"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
)

type Ethtool struct {
	Action   string `json:"action"`
	Link     string `json:"link"`
	Reply    string `json:"reply"`
}

func (req *Ethtool) GetEthTool(rw http.ResponseWriter) {
	_, r := os.Stat(path.Join("/sys/class/net", req.Link))
	if os.IsNotExist(r) {
		log.Errorf("Failed to get link %s: %s", r, req.Link)
		return
	}

	e, r := ethtool.NewEthtool()
	if r != nil {
		log.Errorf("Failed to init ethtool for link %s: %s", r, req.Link)
		return
	}
	defer e.Close()

	switch req.Action {
	case "get-link-stat":
		stats, r := e.Stats(req.Link)
		if r != nil {
			log.Errorf("Failed to get ethtool statitics for link %s: %s", r, req.Link)
			return
		}

		jsonStat, r := json.Marshal(stats)
		if r != nil {
			log.Errorf("Failed to encode ethtool json statitics for link %s: %s", r, req.Link)
			return
		}

		rw.Write([]byte(jsonStat))
		break

	case "get-link-features":

		features, r := e.Features(req.Link)
		if r != nil {
			log.Errorf("Failed to get ethtool features for link %s: %s", r, req.Link)
			return
		}

		jsonFeatures, r := json.Marshal(features)
		if r != nil {
			log.Errorf("Failed to encode json features for link %s: %s", r, req.Link)
			return
		}

		rw.Write([]byte(jsonFeatures))

		break

	case "get-link-bus":

		bus, r := e.BusInfo(req.Link)
		if r != nil {
			log.Errorf("Failed to get ethtool bus for link %s: %s", r, req.Link)
			return
		}

		ethtool := Ethtool {
			Action: "get-link-bus",
			Link:   req.Link,
			Reply:  bus,
		}

		jsonBus, r := json.Marshal(ethtool)
		if r != nil {
			log.Errorf("Failed to get encode json bus information for link %s: %s", r, req.Link)
			return
		}

		rw.Write([]byte(jsonBus))

		break

	case "get-link-driver-name":

		driver, r := e.DriverName(req.Link)
		if r != nil {
			log.Errorf("Failed to get ethtool driver name for link %s: %s", r, req.Link)
			return
		}

		ethtool := Ethtool{
			Action: "get-link-driver-name",
			Link:   req.Link,
			Reply:  driver,
		}

		jsonDriver, r := json.Marshal(ethtool)
		if r != nil {
			log.Errorf("Failed to get encode json driver name for link %s: %s", r, req.Link)
			return
		}

		rw.Write([]byte(jsonDriver))

		break
	}
}
