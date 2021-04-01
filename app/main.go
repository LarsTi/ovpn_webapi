package main

import(
	"log"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)
// User struct for my db
func main(){
	log.Println("Application startup")
	_, err := gorm.Open(sqlite.Open("/docker/data/data.sq3"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err)
	}
}
