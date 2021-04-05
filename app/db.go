package main
import (
	"log"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)
func connDB() (db *DB) {
	db = &DB{}
	conn, err := gorm.Open(sqlite.Open("/docker/data/data.sq3"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err)
	}
	db.conn = conn
	//log.Fatalf raises panic, if there is no db
	return db
}
func (db *DB) init() {
	log.Println("Migrating Structures")
	db.conn.AutoMigrate(&User{})
	db.conn.AutoMigrate(&AccessGroup{})
	db.conn.AutoMigrate(&Certificate{})
	db.conn.AutoMigrate(&UserAccess{})
	log.Println("Migration finished")
}
func (db *DB) loadCA() (ca *CA){
	log.Println("Reading CA from Database")
	ca = &CA{}
	ca.db = db
	ca.ca = db.readCertByCN("root-ca")
	if ca.ca == nil {
		ca.ca = createCA("root-ca")
		ca.SerialOld = ca.ca.Serial
		db.writeCert(ca.ca)
	}else{
		ca.SerialOld = db.getSerialOld()
	}
	return ca
}
func (ca *CA) checkServer(){
	srvCert := ca.db.readCertByCN("server")
	if srvCert != nil {
		log.Println("Found server certificate")
	}else{
		log.Println("Creating server cert")
		srvCert = ca.createServer("server")
		ca.db.writeCert(srvCert)
		log.Println("Created server cert")
	}
	log.Println("Updating key and cert for server (files)")
	srvCert.WriteFileCert()
	srvCert.WriteFileKey()

}
func (db *DB) getSerialOld() (serial int64){
	crt := &Certificate{}
	db.conn.Last(crt)
	if db.conn.Error != nil {
		log.Fatalf("Error: %s", db.conn.Error)
	}
	return crt.Serial
}
func (db *DB) writeCert(cert *Certificate){
	if cert == nil {
		log.Println("Certificate is not bound!")
	}else if cert.CN == "" {
		log.Println("Certificate has no Common Name, discarding from save")
	}else if cert.Private == "" {
		log.Println("Certificate has no private Key, discarding from save")
	}else if cert.Public == "" {
		log.Println("Certificate has no public Key, discarding from save")
	}else if cert.ID == 0{
		log.Printf("Creating entry for cert with common name %s", cert.CN)
		db.conn.Create(cert)
	}else{
		log.Printf("Updating entry for cert with common name %s", cert.CN)
		db.conn.Model(cert).Update("Revoked",cert.Revoked)
	}
}
func (db *DB) readCertByCN(cn string) (cert *Certificate){
	cert = &Certificate{}
	db.conn.First(&cert, "CN = ?", cn)
	if cert.Type == 0 {
		log.Printf("No Certificate for Common Name %s found", cn)
		return nil
	}
	return cert
}
