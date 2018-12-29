// SPDX-License-Identifier: Apache-2.0

package router

import (
	"api-routerd/cmd/share"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	AuthConfPath = "/etc/api-routerd/api-routerd-auth.conf"
)

type TokenDB struct {
	tokenUsers map[string]string
}

func (db *TokenDB) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("X-Session-Token")

		if user, found := db.tokenUsers[token]; found {
			log.Printf("Authenticated user %s\n", user)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
			log.Infof("Unauthorized user")
		}
	})
}

func InitAuthMiddleware()(TokenDB, error) {
	db := TokenDB{make(map[string]string)}

	lines, r := share.ReadFullFile(AuthConfPath)
	if r != nil {
		log.Fatal("Failed to read auth config file")
		return db, errors.New("Failed to read auth config file")
	}

	for _, line := range lines {
		authLine := strings.Fields(line)
		db.tokenUsers[authLine[1]] = authLine[0]
	}

	return db, nil
}
