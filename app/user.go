package main

import (
	"log"
	"net/http"
	"encoding/json"
//	"strconv"
	"os"
	"io"
	"fmt"

	"github.com/gorilla/mux"
)
func (db *DB) loadAllUsers(w http.ResponseWriter, r *http.Request){
	var ret []User
	result := db.conn.Find(&ret)
	log.Printf("Read %d Users from DB\n", result.RowsAffected)
	
	json.NewEncoder(w).Encode(ret)
}
func (db *DB) createUser(w http.ResponseWriter, r *http.Request){
	userIn := User{}
	userDb := User{}
	json.NewDecoder(r.Body).Decode(&userIn)
	db.conn.Where("Mail = ?", userIn.Mail).First(&userDb)
	if userDb.Mail == userIn.Mail {
		log.Printf("User %s exists\n", userIn.Name)
		http.Error(w, "Already Exists", http.StatusBadRequest)
		return
	}
	db.conn.Create(&userIn)
	if (db.conn.Error != nil){
		log.Printf("User Insert error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	userIn.User = userIn.ID
	db.conn.Save(&userIn)
	log.Printf("User %s created with id %d\n", userIn.Name, userIn.ID)
}
func (db *DB) updateUser(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	userIn := User{}
	userDb := User{}
	json.NewDecoder(r.Body).Decode(&userIn)
	
	result := db.conn.Where("ID = ?", id).First(&userDb)
	if (db.conn.Error != nil){
		log.Printf("User Update error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}else if (result.RowsAffected == 0){
		log.Printf("User Update error (user not found)\n")
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	userDb.Name = userIn.Name
	userDb.Surname = userIn.Surname
	userDb.Org = userIn.Org
	userDb.Mail = userIn.Mail
	userDb.Passwd = userIn.Passwd
	db.conn.Save(&userDb)
	if (db.conn.Error != nil){
		log.Printf("User Update error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("User %s updated with id %d\n", userIn.Name, id)
}
func (db *DB) deleteUser(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	db.conn.Delete(&User{}, id)
	if (db.conn.Error != nil){
		log.Printf("User Delete error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	
	for _, cert := range db.getCertsForUser(id){
		log.Printf("Revoking cert %s for user %s\n", cert.CN, id)
		db.revokeCert(&cert)
	}
	db.loadCA().createCRL(db.getRevokedCerts())

	db.createCCD(id)
	log.Printf("User with id %s deleted\n", id)
}
func (db *DB) getCertsForUser(userId string) (ret []Certificate){
	result := db.conn.Where("mail = ?", userId).Find(&ret)
	log.Printf("Found %d active Certificates for user %s\n", result.RowsAffected, userId)
	return ret
}
func (db *DB) writePWFile(){
	var users []User
	db.conn.Find(&users)
	file, err := os.OpenFile("/docker/server/pw", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if (err != nil){
		log.Println(err)
		return
	}
	defer file.Close()
	for _, user := range users{
		_, err = io.WriteString(file,fmt.Sprintf("%s.%s:%s", user.Surname, user.Name, user.Passwd))
		if (err != nil) {
			log.Printf("Could not write pw file for user %s.%s", user.Surname, user.Name)
			log.Println(err)
		}else{
			log.Printf("Wrote pw file for user %s.%s", user.Surname, user.Name)
		}
	}
	log.Println("Finished write of pw file")
}
