// SPDX-License-Identifier: Apache-2.0

package network

import (
	"api-routerd/cmd/network/ethtool"
	"api-routerd/cmd/network/networkd"
	"api-routerd/cmd/network/resolve"
	"api-routerd/cmd/network/systemdresolved"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func NetworkLinkGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	l := Link {Link: link}

	switch r.Method {
	case "GET":
		l.GetLink(rw)
		break
	}
}

func NetworkLinkAdd(rw http.ResponseWriter, r *http.Request) {
	link := new(Link)

	err := json.NewDecoder(r.Body).Decode(&link);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		switch link.Action {
		case "add-link-bridge":
			err := link.LinkCreateBridge()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func NetworkLinkDelete(rw http.ResponseWriter, r *http.Request) {
	link := new(Link)

	err := json.NewDecoder(r.Body).Decode(&link);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "DELETE":
		switch link.Action {
		case "delete-link":
			err := link.LinkDelete()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func NetworkLinkSet(rw http.ResponseWriter, r *http.Request) {
	link := new(Link)

	err := json.NewDecoder(r.Body).Decode(&link);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "PUT":
		switch link.Action {
		case "set-link-up", "set-link-down", "set-link-mtu":
			err := link.SetLink()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func NetworkGetAddress(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	address := Address{Link: link}

	switch r.Method {
	case "GET":
		address.GetAddress(rw)
		break
	}
}

func NetworkAddAddress(rw http.ResponseWriter, r *http.Request) {
	address := new(Address)

	err := json.NewDecoder(r.Body).Decode(&address);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		switch address.Action {
		case "add-address":
			err := address.AddAddress()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		break
	}
}

func NetworkDeleteAddres(rw http.ResponseWriter, r *http.Request) {
	address := new(Address)

	err := json.NewDecoder(r.Body).Decode(&address);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "DELETE":
		err := address.DelAddress()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	}
}

func NetworkAddRoute(rw http.ResponseWriter, r *http.Request) {
	route := new(Route)

	err := json.NewDecoder(r.Body).Decode(&route);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		switch route.Action {
		case "add-default-gw":
			err := route.AddDefaultGateWay()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			break
		}

		break
	case "PUT":
		switch route.Action {
		case "replace-default-gw":
			err := route.ReplaceDefaultGateWay()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			break
		}
	}
}

func NetworkDeleteRoute(rw http.ResponseWriter, r *http.Request) {
	route := new(Route)

	err := json.NewDecoder(r.Body).Decode(&route);
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "DELETE":
		switch route.Action {
		case "del-default-gw":
			err := route.DelDefaultGateWay()
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			break
		}
	}
}

func NetworkGetRoute(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetRoutes(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func NetworkdConfigureNetwork(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		networkd.ConfigureNetworkFile(rw, r)
		break
	}
}

func NetworkdConfigureNetDev(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		networkd.ConfigureNetDevFile(rw, r)
		break
	}
}

func NetworkdConfigureLink(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		networkd.ConfigureLinkFile(rw, r)
		break
	}
}

func NetworkConfigureEthtool(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]
	command := vars["command"]

	ethtool := ethtool.Ethtool{Link: link, Action: command}

	switch r.Method {
	case "GET":

		err := ethtool.GetEthTool(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func NetworkConfigureResolv(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := resolv.GetResolvConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := resolv.UpdateResolvConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := resolv.DeleteResolvConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func NetworkConfigureSystemdResolved(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := systemdresolved.GetResolveConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := systemdresolved.UpdateResolveConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := systemdresolved.DeleteResolveConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

func RegisterRouterNetwork(router *mux.Router) {
	n := router.PathPrefix("/network").Subrouter()

	// Link
	n.HandleFunc("/link/set",        NetworkLinkSet)
	n.HandleFunc("/link/add",        NetworkLinkAdd)
	n.HandleFunc("/link/delete",     NetworkLinkDelete)
	n.HandleFunc("/link/get/{link}", NetworkLinkGet)
	n.HandleFunc("/link/get",        NetworkLinkGet)

	// Address
	n.HandleFunc("/address/add",        NetworkAddAddress)
	n.HandleFunc("/address/delete",     NetworkDeleteAddres)
	n.HandleFunc("/address/get",        NetworkGetAddress)
	n.HandleFunc("/address/get/{link}", NetworkGetAddress)

	// Route
	n.HandleFunc("/route/add", NetworkAddRoute)
	n.HandleFunc("/route/del", NetworkDeleteRoute)
	n.HandleFunc("/route/get", NetworkGetRoute)

	// systemd-networkd
	networkd.InitNetworkd()
	n.HandleFunc("/networkd/network", NetworkdConfigureNetwork)
	n.HandleFunc("/networkd/netdev",  NetworkdConfigureNetDev)
	n.HandleFunc("/networkd/link",    NetworkdConfigureLink)

	// ethtool
	n.HandleFunc("/ethtool/{link}/{command}", NetworkConfigureEthtool)

	// resolv.conf
	n.HandleFunc("/resolv", NetworkConfigureResolv)
	n.HandleFunc("/resolv/get", NetworkConfigureResolv)
	n.HandleFunc("/resolv/add", NetworkConfigureResolv)
	n.HandleFunc("/resolv/delete", NetworkConfigureResolv)

	// systemd-resolved
	n.HandleFunc("/systemdresolved", NetworkConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/get", NetworkConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/add", NetworkConfigureSystemdResolved)
	n.HandleFunc("/systemdresolved/delete", NetworkConfigureSystemdResolved)

}
