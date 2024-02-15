package main

import(
	"log"
	"os"

)
// User struct for my db
func main(){
	log.Println("Application startup")
	log.Println("Please make sure, /docker/data, /docker/ccd and optional /docker/server is writable for execution user!")
	log.Println("Checking for prototype Files")
	checkFile("/docker/data/client.ovpn.proto")
	log.Println("Finished check. Any Errors? It will not hardly fail, but i doubt it will work as intended")

	db := getSingleton().dbConn
	
	//db := connDB()
	//db.init()
	ca := db.loadCA()
	
	getSingleton().ca.WriteFileCert()
	
	ca.checkServer()
	ca.createCRL(db.getRevokedCerts())
	
	getSingleton().dbConn.writePWFile()

	//this one blocks!
	RunGin(8080, db)
	//RunWebApi(8080, db)
}
func checkFile(file string){
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		log.Printf("The file %s seems not to exist. Sure the program will work?!\n", file)
		return
	}
	log.Printf("Found file %s. Assuming it is correct.\n", file)
}
