package main

import (
	"log"
	"net/http"
	"encoding/json"
//	"strconv"

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
	userIn.UserId = userIn.ID
	db.conn.Save(&userIn)
	log.Printf("User %s created with id %d\n", userIn.Name, userIn.ID)
}
func (db *DB) updateUser(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	userIn := User{}
	userDb := User{}
	json.NewDecoder(r.Body).Decode(&userIn)
	
	db.conn.Where("ID = ?", id).First(&userDb)
	if (db.conn.Error != nil){
		log.Printf("User Update error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	userDb.Name = userIn.Name
	userDb.Surname = userIn.Surname
	userDb.Org = userIn.Org
	userDb.Mail = userIn.Mail

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
	//TODO: revoke Certs and correct ccd
	log.Printf("User with id %s deleted\n", id)
}
