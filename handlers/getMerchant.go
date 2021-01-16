package handlers

import (
	"go-merchants/database"
	"net/http"

	"github.com/gorilla/mux"
)

// GetMerchant handler for getting existing merchants details in the database
func GetMerchant(db database.MerchantsInterface, membersDb database.MembersInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		res, err := db.Get(id)
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

		res.MembersCount, _ = membersDb.GetMembersCountByMerchantID(res.ID)

		WriteResponse(w, http.StatusOK, res)
	}
}
