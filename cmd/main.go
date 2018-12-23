// SPDX-License-Identifier: Apache-2.0

package main

import (
	"api-routerd/cmd/router"
	"api-routerd/cmd/share"
)

func main() {
	share.InitLog()
	router.StartRouter()
}
