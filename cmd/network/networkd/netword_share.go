// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"api-routerd/cmd/share"
	log "github.com/sirupsen/logrus"
)

const NetworkdUnitPath = "/var/run/systemd/network"

func InitNetworkd() (err error) {
	r := share.CreateDirectory(NetworkdUnitPath, 0777)
	if (r != nil) {
		log.Errorf("Failed create network unit path %s: %s", NetworkdUnitPath, r)
		return r
	}

	return nil
}
