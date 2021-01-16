package handlers

import (
	"go-merchants/database"
	"net/http"

	"github.com/gorilla/mux"
)

// DeleteMember handler for deleting existing members in the database
func DeleteMember(db database.MembersInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		res, err := db.Delete(id)

		if err != nil {
			if err.Error() == "No data" {
				WriteResponse(w, http.StatusBadRequest, map[string]string{
					"message": "Member doesn't exist",
				})
				return
			}

			WriteResponse(w, http.StatusBadGateway, err.Error())
			return
		}

		WriteResponse(w, http.StatusOK, res)
	}
}
