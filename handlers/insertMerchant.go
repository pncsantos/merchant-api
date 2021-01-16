package handlers

import (
	"encoding/json"
	"go-merchants/database"
	"go-merchants/models"
	"io/ioutil"
	"net/http"
)

// InsertMerchant handler for adding new merchants in the database
func InsertMerchant(db database.MerchantsInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchant := models.Merchant{}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &merchant)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := db.Insert(merchant)
		if err != nil {
			WriteResponse(w, http.StatusBadGateway, err.Error())
			return
		}

		WriteResponse(w, http.StatusOK, res)
	}
}
