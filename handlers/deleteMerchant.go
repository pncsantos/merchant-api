package handlers

import (
	"go-merchants/database"
	"net/http"

	"github.com/gorilla/mux"
)

// DeleteMerchant handler for deleting existing merchants in the database
func DeleteMerchant(db database.MerchantsInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		res, err := db.Delete(id)
		if err != nil {
			if err.Error() == "No records found" {
				WriteResponse(w, http.StatusBadRequest, map[string]string{
					"message": "Merchant doesn't exist",
				})
				return
			}

			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		WriteResponse(w, http.StatusOK, res)
	}
}
