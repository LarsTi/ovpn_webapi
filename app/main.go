package main

import(
	"log"
)
// User struct for my db
func main(){
	log.Println("Application startup")

	db := connDB()
	db.init()
	ca := db.loadCA()
	ca.ca.WriteFileCert()
	//ca.ca.WriteCRL()
	ca.checkServer()
	ca.createCRL(db.getRevokedCerts())
}
