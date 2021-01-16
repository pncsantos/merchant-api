package handlers

import (
	"encoding/json"
	"go-merchants/database"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// UpdateMerchant handler for updating merchant details by ID in the database
func UpdateMerchant(db database.MerchantsInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var merchant interface{}
		err = json.Unmarshal(body, &merchant)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := db.Update(id, merchant)
		if err != nil {
			if err.Error() == "No records found" {
				WriteResponse(w, http.StatusBadRequest, map[string]string{
					"message": "Merchant doesn't exist",
				})
				return
			}
			WriteResponse(w, http.StatusBadGateway, err.Error())
			return
		}

		WriteResponse(w, http.StatusOK, res)
	}
}

// WriteResponse update response by status
func WriteResponse(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
