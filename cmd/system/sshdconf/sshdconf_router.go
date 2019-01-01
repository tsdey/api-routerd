// SPDX-License-Identifier: Apache-2.0

package sshdconf

import (
        "encoding/json"
        "github.com/gorilla/mux"
        "net/http"
)

//function to read sudoers rule
func ReadConfigSudoers(rw http.ResponseWriter, req *http.Request) {
        sshdconf := new(SSHdConf)
        _ = json.NewDecoder(req.Body).Decode(&sshdconf)
        GetSudoers(rw)
}

//function check if sshd configuration exists
func ReadSSHConfigAvailable(rw http.ResponseWriter, req *http.Request) {
        sshdconf := new(SSHdConf)
        _ = json.NewDecoder(req.Body).Decode(&sshdconf)
        SSHConfigAvailable(rw)
}

//function to read sshd configuration file
func ReadSSHConfigRead(rw http.ResponseWriter, req *http.Request) {
        sshdconf := new(SSHdConf)
        _ = json.NewDecoder(req.Body).Decode(&sshdconf)
        SSHConfFileRead(rw)
}

//register router to read configuration files
//for various services
func RegisterRouterSSHdConf(router *mux.Router) {
        s := router.PathPrefix("/system").Subrouter().StrictSlash(false)
        s.HandleFunc("/sshd", ReadSSHConfigAvailable)
        s.HandleFunc("/sshd/sudoers", ReadConfigSudoers)
        s.HandleFunc("/sshd/sshdconf", ReadSSHConfigRead)
}

