package main

import (
	"log"
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
)
func (db *DB) loadUserCertificates(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	ret := db.getCertsForUser(id)
	
	json.NewEncoder(w).Encode(ret)
}
func (db *DB) createUserCertificate(w http.ResponseWriter, r *http.Request){
	certIn := Certificate{}
	userDb := User{}
	id := mux.Vars(r)["id"]

	json.NewDecoder(r.Body).Decode(&certIn)
	idString := fmt.Sprintf("%d", certIn.User)
	if (idString != id) {
		log.Printf("User (%s) and ID (%s) field not equal, abort!\n", idString, id)
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	result := db.conn.Where("ID = ?", id).First(&userDb)
	if (db.conn.Error != nil){
		log.Printf("Certificate Create error (user not found): %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}else if (result.RowsAffected == 0){
		log.Printf("Certificate Create error (user not found)\n")
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	
	ca := db.loadCA()
	crt := ca.createClient(&userDb)
	
	db.conn.Create(&crt)

	db.createCCD(id)

	log.Printf("Certificate %s created for user %d\n", crt.CN, crt.User)
}
func (db *DB) deleteUserCertificate(w http.ResponseWriter, r *http.Request){
	crt := Certificate{}
	id := mux.Vars(r)["id"]
	cert := mux.Vars(r)["cert"]
	log.Printf("Revoking Certificate %s for user %s\n", cert, id)
	
	result := db.conn.Where("User = $1 AND ID = $2", id, cert).Delete(&crt)
	if (db.conn.Error != nil){
		log.Printf("Certificate Delete error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}else if (result.RowsAffected == 0){
		log.Printf("Certification revokation error (user/certificate not found)\n")
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	db.revokeCert(&crt)

	db.loadCA().createCRL(db.getRevokedCerts())

	db.createCCD(id)

	log.Printf("Certificate for user %s with ID %s revoked\n", id, cert)
}
