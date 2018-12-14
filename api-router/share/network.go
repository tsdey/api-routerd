// SPDX-License-Identifier: Apache-2.0

package share

import (
	"strings"
	"golang.org/x/sys/unix"
)

func IsValidIfName(ifname string) (bool) {
	s := strings.TrimSpace(ifname)
	if (len(s) == 0 || len(s) > unix.IFNAMSIZ) {
		return false
	}

	return true
}
