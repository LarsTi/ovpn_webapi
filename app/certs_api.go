package main

import (
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"strings"
	"io/ioutil"

	"github.com/gorilla/mux"
)
func (db *DB) loadUserCertificates(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	ret := db.getCertsForUser(id)
	
	json.NewEncoder(w).Encode(ret)
}
func (db *DB) createUserCertificate(w http.ResponseWriter, r *http.Request){
	userDb := User{}
	id := mux.Vars(r)["id"]

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
	
	result := db.conn.Where("User = ? AND ID = ?", id, cert).Delete(&crt)
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
func (db *DB) downloadUserCertificate(w http.ResponseWriter, r *http.Request){
	crt := Certificate{}
	id := mux.Vars(r)["id"]
	cert := mux.Vars(r)["cert"]
	log.Printf("creating config for for user %s and cert %s\n", id, cert)
	
	result := db.conn.Where("User = ? AND ID = ?", id, cert).First(&crt)
	if (db.conn.Error != nil){
		log.Printf("Certificate Read error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}else if (result.RowsAffected == 0){
		log.Printf("Certification read error (user/certificate not found)\n")
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	
	ca := db.loadCA()
	content, err := ioutil.ReadFile("/docker/data/client.ovpn.proto")
	if err != nil {
		log.Printf("Error reading proto file: %s", err)
	}
	lines := strings.Split(string(content), "\n")
	lines = append(lines, "<ca>")
	lines = append(lines, ca.ca.Public)
	lines = append(lines, "</ca>")

	lines = append(lines, "<key>")
	lines = append(lines, crt.Private)
	lines = append(lines, "</key>")

	lines = append(lines, "<cert>")
	lines = append(lines, crt.Public)
	lines = append(lines, "</cert>")


	w.Header().Set("Content-Type", "text/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s.ovpn\"", crt.CN))
	fmt.Fprintf(w, strings.Join(lines, "\n"))

	db.createCCD(id)

	log.Printf("Certificate for user %s with ID %s build and returned\n", id, cert)
}
