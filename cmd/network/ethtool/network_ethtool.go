// SPDX-License-Identifier: Apache-2.0

package ethtool

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"errors"
	"github.com/safchain/ethtool"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Ethtool struct {
	Action   string `json:"action"`
	Link     string `json:"link"`
}

func (req *Ethtool) GetEthTool(rw http.ResponseWriter) (error) {
	var j []byte

	link := share.LinkExists(req.Link)
	if !link {
		log.Errorf("Failed to get link: %s", req.Link)
		return errors.New("Link not found")
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		log.Errorf("Failed to init ethtool for link %s: %s", err, req.Link)
		return err
	}
	defer e.Close()

	switch req.Action {
	case "get-link-stat":
		stats, err := e.Stats(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool statitics for link %s: %s", err, req.Link)
			return err
		}

		j, err = json.Marshal(stats)
		if err != nil {
			log.Errorf("Failed to encode ethtool json statitics for link %s: %s", err, req.Link)
			return err
		}
		break

	case "get-link-features":

		features, err := e.Features(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool features for link %s: %s", err, req.Link)
			return err
		}

		j, err = json.Marshal(features)
		if err != nil {
			log.Errorf("Failed to encode json features for link %s: %s", err, req.Link)
			return err
		}
		break

	case "get-link-bus":

		bus, err := e.BusInfo(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool bus for link %s: %s", err, req.Link)
			return err
		}

		b := struct {
			Bus string
		}{
			bus,
		}

		j, err = json.Marshal(b)
		if err != nil {
			log.Errorf("Failed to get encode json bus information for link %s: %s", err, req.Link)
			return err
		}

		break

	case "get-link-driver-name":

		driver, err := e.DriverName(req.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool driver name for link %s: %s", err, req.Link)
			return err
		}

		d := struct {
			Driver string
		}{
			driver,
		}

		j, err = json.Marshal(d)
		if err != nil {
			log.Errorf("Failed to get encode json driver name for link %s: %s", err, req.Link)
			return err
		}

		break
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(j)

	return nil
}
