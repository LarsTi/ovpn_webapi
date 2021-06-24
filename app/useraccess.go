package main

import (
	"log"
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
)
func (db *DB) loadUserAccess(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	var ret []UserAccess
	result := db.conn.Where("UserId = ?", id).Find(&ret)
	log.Printf("Read %d UserAccess from DB for UserId %s\n", result.RowsAffected, id)
	
	json.NewEncoder(w).Encode(ret)
}
func (db *DB) createUserAccess(w http.ResponseWriter, r *http.Request){
	userIn := UserAccess{}
	userDb := User{}
	agDb := AccessGroup{}
	id := mux.Vars(r)["id"]

	json.NewDecoder(r.Body).Decode(&userIn)
	idString := fmt.Sprintf("%d", userIn.UserId)
	if (idString != id) {
		log.Printf("UserId (%s) and ID (%s) field not equal, abort!\n", idString, id)
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	db.conn.Where("ID = ?", id).First(&userDb)
	if (db.conn.Error != nil){
		log.Printf("UserAccess Create error (user not found): %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	db.conn.Where("ID = ?", userIn.AccessGroup).First(&agDb)
	if (db.conn.Error != nil){
		log.Printf("UserAccess Create error (accessgroup not found): %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	
	db.conn.Create(&userIn)

	db.createCCD(id)

	log.Printf("UserAccess %s created with id %d\n", userIn.AccessGroup, userIn.ID)
}
func (db *DB) deleteUserAccess(w http.ResponseWriter, r *http.Request){
	ua := UserAccess{}
	id := mux.Vars(r)["id"]
	group := mux.Vars(r)["group"]

	db.conn.Where("UserId = ?", id).Where("AccessGroup = ?", group).Delete(&ua)
	if (db.conn.Error != nil){
		log.Printf("UserAccess Delete error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	
	db.createCCD(id)

	log.Printf("UserAccess with for user %s and group %s deleted\n", id, group)
}
func (db *DB) createCCD(userId string){
	log.Printf("Creating CCD files for User %s\n", userId)
}

