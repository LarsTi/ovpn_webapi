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
	SerialOld	int64		//letzte Vergebene Serial
}

var singleInstance *singleton


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
				singleInstance.SerialOld = singleInstance.ca.Serial
				singleInstance.dbConn.writeCert(singleInstance.ca)
			}else{
				singleInstance.SerialOld = singleInstance.dbConn.getSerialOld()
			}

        } else {
            //fmt.Println("Single instance already created.")
        }
    } else {
        //fmt.Println("Single instance already created.")
    }

    return singleInstance
}
