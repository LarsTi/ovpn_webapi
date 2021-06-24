package main

import (
	"log"
	"net/http"
	"encoding/json"
//	"strconv"

	"github.com/gorilla/mux"
)
func (db *DB) loadAllAccessGroups(w http.ResponseWriter, r *http.Request){
	var ret []AccessGroup
	result := db.conn.Find(&ret)
	log.Printf("Read %d AccessGroups from DB\n", result.RowsAffected)
	
	json.NewEncoder(w).Encode(ret)
}
func (db *DB) createAccessGroup(w http.ResponseWriter, r *http.Request){
	agIn := AccessGroup{}
	agDb := AccessGroup{}
	json.NewDecoder(r.Body).Decode(&agIn)
	db.conn.Where("name = ?", agIn.Name).First(&agDb)
	if agDb.Name == agIn.Name {
		log.Printf("AccessGroup %s exists\n", agIn.Name)
		http.Error(w, "Already Exists", http.StatusBadRequest)
		return
	}
	db.conn.Create(&agIn)
	if (db.conn.Error != nil){
		log.Printf("AccessGroup Insert error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("AccessGroup %s created with id %d\n", agIn.Name, agIn.ID)
}
func (db *DB) updateAccessGroup(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	agIn := AccessGroup{}
	agDb := AccessGroup{}
	json.NewDecoder(r.Body).Decode(&agIn)
	
	db.conn.Where("ID = ?", id).First(&agDb)
	if (db.conn.Error != nil){
		log.Printf("AccessGroup Update error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	agDb.Name = agIn.Name
	agDb.Subnet = agIn.Subnet
	agDb.Mask = agIn.Mask
	db.conn.Save(&agDb)
	if (db.conn.Error != nil){
		log.Printf("AccessGroup Update error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("AccessGroup %s updated with id %d\n", agIn.Name, id)
}
func (db *DB) deleteAccessGroup(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	db.conn.Delete(&AccessGroup{}, id)
	if (db.conn.Error != nil){
		log.Printf("AccessGroup Delete error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("AccessGroup with id %s deleted\n", id)
}

