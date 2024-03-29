package main

import (
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	"io"
	"strconv"

	"github.com/gorilla/mux"
)
func (db *DB) loadUserAccess(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	var ret []UserAccess
	result := db.conn.Where("User = ?", id).Find(&ret)
	log.Printf("Read %d UserAccess from DB for User %s\n", result.RowsAffected, id)
	
	json.NewEncoder(w).Encode(ret)
}
func (db *DB) createUserAccess(w http.ResponseWriter, r *http.Request){
	userIn := UserAccess{}
	uaDb := UserAccess{}
	userDb := User{}
	agDb := AccessGroup{}
	id := mux.Vars(r)["id"]

	json.NewDecoder(r.Body).Decode(&userIn)
	idString := fmt.Sprintf("%d", userIn.User)
	if (idString != id) {
		log.Printf("User (%s) and ID (%s) field not equal, abort!\n", idString, id)
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	result := db.conn.Where("ID = ?", id).First(&userDb)
	if (db.conn.Error != nil){
		log.Printf("UserAccess Create error (user not found): %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}else if (result.RowsAffected == 0){
		log.Printf("UserAccess Create error (user not found)\n")
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	
	result = db.conn.Where("ID = ?", userIn.Access).First(&agDb)
	if (db.conn.Error != nil){
		log.Printf("UserAccess Create error (accessgroup not found): %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}else if (result.RowsAffected == 0){
		log.Printf("UserAccess Create error (group not found)\n")
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	
	result = db.conn.Where("User = ? AND Access = ?", id, agDb.ID).Find(&uaDb)
	if (db.conn.Error != nil){
		log.Printf("UserAccess Create error (syntax): %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}else if (result.RowsAffected > 0){
		log.Printf("UserAccess Create error (group already added)\n")
		http.Error(w, "Exists", http.StatusBadRequest)
		return
	}

	db.conn.Create(&userIn)

	db.createCCD(id)

	log.Printf("UserAccess %d created with id %d\n", userIn.Access, userIn.ID)
}
func (db *DB) deleteUserAccess(w http.ResponseWriter, r *http.Request){
	ua := UserAccess{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	group, _ := strconv.Atoi(mux.Vars(r)["group"])
	log.Printf("Deleting AccessGroup %s for user %s\n", group, id)
	
	result := db.conn.Where("User = ?", id).Where("Access = ?", group).Delete(&ua)
	log.Println(ua)
	log.Printf("id: %T\n", id)
	log.Printf("group: %T\n", group)
	if (db.conn.Error != nil){
		log.Printf("UserAccess Delete error: %s\n", db.conn.Error)
		http.Error(w, db.conn.Error.Error(), http.StatusBadRequest)
		return
	}else if (result.RowsAffected == 0){
		log.Printf("UserAccess Delete error (user/group not found)\n")
		http.Error(w, "Wrong call", http.StatusBadRequest)
		return
	}
	
	db.createCCD(fmt.Sprintf("%d",id))

	log.Printf("UserAccess with for user %s and group %s deleted\n", id, group)
}
func (db *DB) createCCD(mail string){
	log.Printf("Correcting CCD files for User %s\n", mail)
	
	certs := getSingleton().dbConn.getCertsForUser(mail)
	for _, c := range certs {
		writeCCD(c)
	}
}
func writeCCD(c Certificate){
	path := fmt.Sprintf("/docker/ccd/%s", c.CN)
	log.Printf("Writing file: %s\n", path)
	file, err := os.Create(path)
	if (err != nil){
		log.Println(err)
		return
	}
	defer file.Close()
	for _, ag := range getSingleton().dbConn.getAccessGroupForUser(c.Mail) {
		_, err = io.WriteString(file, fmt.Sprintf("\n#%s\npush \"route %s %s\"\n", ag.Name, ag.Subnet, ag.Mask))
		if (err != nil){
			log.Println(err)
			continue
		}
	}
	log.Printf("Finished file: %s\n", path)

}
