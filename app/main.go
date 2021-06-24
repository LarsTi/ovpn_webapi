package main

import(
	"log"
)
// User struct for my db
func main(){
	log.Println("Application startup")
	log.Println("Please make sure, /docker/data is writable for execution user!")
	db := connDB()
	db.init()
	ca := db.loadCA()
	ca.ca.WriteFileCert()
	
	ca.checkServer()
	ca.createCRL(db.getRevokedCerts())

	//this one blocks!
	RunWebApi(8080, db)
}
