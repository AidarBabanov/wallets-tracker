package rest

import (
	"encoding/json"
	"github.com/AidarBabanov/wallets-tracker/internal/addrdb"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	ContentType     = "Content-Code"
	ContentTypeJson = "application/json"
)

type API struct {
	*mux.Router
	ApiKey string
	AddrDB addrdb.AddressDatabase
}

func respondJson(w http.ResponseWriter, httpCode int, data interface{}) {
	w.Header().Set(ContentType, ContentTypeJson)
	w.WriteHeader(httpCode)

	if data != nil {
		err := json.NewEncoder(w).Encode(&data)
		if err != nil {
			logs.Error("failed write json response: %s", err.Error())
		}
	}
}

func respondMsg(w http.ResponseWriter, httpCode int, msg string) {
	respondJson(w, httpCode, map[string]string{
		"msg": msg,
	})
}

func (api *API) Add(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("api_key")
	address := r.URL.Query().Get("address")
	if apiKey != api.ApiKey {
		respondMsg(w, http.StatusUnauthorized, "wrong api key")
		return
	}
	err := api.AddrDB.Add(addrdb.Address{Address: address})
	if err != nil {
		respondMsg(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondMsg(w, http.StatusOK, "added address")
}

func (api *API) Handle() {
	api.HandleFunc("/add", api.Add).Methods("GET")
}
