package handlers

import (
	"go-merchants/database"
	"net/http"

	"github.com/gorilla/mux"
)

// GetMember handler for getting existing members details in the database
func GetMember(db database.MembersInterface) http.HandlerFunc {
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

		WriteResponse(w, http.StatusOK, res)
	}
}
