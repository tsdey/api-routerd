// SPDX-License-Identifier: Apache-2.0

package share

import (
	"github.com/godbus/dbus"
)

func GetSystemBusPrivateConn() (conn *dbus.Conn, err error) {
	conn, err = dbus.SystemBusPrivate()
	if err != nil {
		return nil, err
	}

	if err = conn.Auth(nil); err != nil {
		conn.Close()
		conn = nil
		return
	}

	if err = conn.Hello(); err != nil {
		conn.Close()
		conn = nil
	}

	return conn, nil
}
