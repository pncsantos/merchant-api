package main

import (
	"context"
	"go-merchants/config"
	"go-merchants/database"
	"go-merchants/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	conf := config.GetConfig()
	ctx := context.TODO()

	db := database.ConnectDB(ctx, conf.Mongo)

	merchantCollections := db.Collection(conf.Mongo.MerchantsCollection)
	merchantClient := &database.MerchantClient{
		Col: merchantCollections,
		Ctx: ctx,
	}

	memberCollections := db.Collection(conf.Mongo.MembersCollection)
	memberClient := &database.MemberClient{
		Col: memberCollections,
		Ctx: ctx,
	}

	r := mux.NewRouter()

	// merchants
	r.HandleFunc("/merchants", handlers.SearchMerchants(merchantClient, memberClient)).Methods("GET")
	r.HandleFunc("/merchants/{id}", handlers.GetMerchant(merchantClient, memberClient)).Methods("GET")
	r.HandleFunc("/merchants", handlers.InsertMerchant(merchantClient)).Methods("POST")
	r.HandleFunc("/merchants/{id}", handlers.UpdateMerchant(merchantClient)).Methods("PATCH")
	r.HandleFunc("/merchants/{id}", handlers.DeleteMerchant(merchantClient)).Methods("DELETE")
	r.HandleFunc("/merchants/{id}/upload", handlers.UploadMerchantPhoto(merchantClient)).Methods("PUT")

	// members
	r.HandleFunc("/merchants/{id}/members", handlers.SearchMembers(memberClient)).Methods("GET")
	r.HandleFunc("/members/{id}", handlers.GetMember(memberClient)).Methods("GET")
	r.HandleFunc("/members", handlers.InsertMember(memberClient)).Methods("POST")
	r.HandleFunc("/members/{id}", handlers.UpdateMember(memberClient)).Methods("PATCH")
	r.HandleFunc("/members/{id}", handlers.DeleteMember(memberClient)).Methods("DELETE")

	// CORS setup
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "PUT"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
