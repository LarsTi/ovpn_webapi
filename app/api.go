package main

import(
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
)

func RunWebApi(port int, db *DB){
	router := mux.NewRouter()
	
	router.HandleFunc("/api/accessgroup", db.loadAllAccessGroups).Methods("GET")
	router.HandleFunc("/api/accessgroup", db.createAccessGroup).Methods("POST")
	router.HandleFunc("/api/accessgroup/{id}", db.deleteAccessGroup).Methods("DELETE")
	router.HandleFunc("/api/accessgroup/{id}", db.updateAccessGroup).Methods("UPDATE")

	router.HandleFunc("/api/user", db.loadAllUsers).Methods("GET")
	router.HandleFunc("/api/user", db.createUser).Methods("POST")
	router.HandleFunc("/api/user/{id}", db.deleteUser).Methods("DELETE")
	router.HandleFunc("/api/user/{id}", db.updateUser).Methods("UPDATE")

	router.HandleFunc("/api/user/{id}/access", db.loadUserAccess).Methods("GET")
	router.HandleFunc("/api/user/{id}/access", db.createUserAccess).Methods("POST")
	router.HandleFunc("/api/user/{id}/access/{group}", db.deleteUserAccess).Methods("DELETE")

	router.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("[webapi-request] %s: %s\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
