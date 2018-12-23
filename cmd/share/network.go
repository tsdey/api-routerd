// SPDX-License-Identifier: Apache-2.0

package share

import (
	"strings"
	"golang.org/x/sys/unix"
	"path"
	"os"
)

func IsValidIfName(ifname string) (bool) {
	s := strings.TrimSpace(ifname)
	if (len(s) == 0 || len(s) > unix.IFNAMSIZ) {
		return false
	}

	return true
}

func LinkExists(ifname string) (bool) {
	_, r := os.Stat(path.Join("/sys/class/net", ifname))
	if os.IsNotExist(r) {
		return false
	}

	return true
}
