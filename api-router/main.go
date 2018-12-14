// SPDX-License-Identifier: Apache-2.0

package main

import (
	"restgateway/api-router/router"
	"restgateway/api-router/share"
)

func main() {
	share.InitLog()
	router.StartRouter()
}
