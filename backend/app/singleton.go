package main
import (
    "fmt"
    "sync"
)

var lock = &sync.Mutex{}

type singleton struct {
	dbConn		*DB
	dbLock 		*sync.Mutex
	ca			*Certificate	//CA
}

var singleInstance *singleton

func (singleton *singleton) getSerialOld()(serialOld int64){
	return singleton.getDb().getSerialOld()
}
func (singleton *singleton) getDb()(db *DB){
	singleton.dbLock.Lock()
	defer singleton.dbLock.Unlock()

	return singleton.dbConn

}
func getSingleton() *singleton {
    if singleInstance == nil {
		lock.Lock()
        defer lock.Unlock()
        if singleInstance == nil {
            fmt.Println("Creating single instance now.")
            singleInstance = &singleton{}
			singleInstance.dbLock = lock
			singleInstance.dbConn = connDB()
			singleInstance.dbConn.init()
			singleInstance.ca = singleInstance.dbConn.readCertByCN("ca")
			if singleInstance.ca == nil {
				singleInstance.ca = createCA("ca")
				singleInstance.dbConn.writeCert(singleInstance.ca)
			}

        } else {
            //fmt.Println("Single instance already created.")
        }
    } else {
        //fmt.Println("Single instance already created.")
    }

    return singleInstance
}
