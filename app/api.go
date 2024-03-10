package main

import(
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)
func RunGin(port int, db *DB){
	config := ginkeycloak.BuilderConfig{
		Service: "vpn-frontend",
		Url: "https://sso.tilsner.io",
		Realm: "master",
	}
	r := gin.Default()
	
	endpointCerts(r, config)
	endpointUser(r, config)
	endpointAccessDefinition(r, config)
	
	log.Println("Starting normal API operation!")
	r.Run()
}
func endpointUser(router *gin.Engine, config ginkeycloak.BuilderConfig){
	secured_realm := router.Group("/api/user")
	secured_realm.Use(ginkeycloak.NewAccessBuilder(config).
		RestrictButForRole("vpn-admin-admin").
		Build())
}
func endpointAccessDefinition(router *gin.Engine, config ginkeycloak.BuilderConfig){
	
}
func endpointCerts(router *gin.Engine, config ginkeycloak.BuilderConfig){
	secured_realm := router.Group("/api/certs")
	secured_realm.Use(ginkeycloak.NewAccessBuilder(config).
		RestrictButForRole("vpn-admin-user").
		Build())

		//Alle Zertifikate anzeigen für den Benutzer
	secured_realm.GET("/", func(c *gin.Context) {
		ginToken, _ := c.Get("token")
		token := ginToken.(ginkeycloak.KeyCloakToken)
		c.JSON(http.StatusOK, getSingleton().dbConn.getCertsForUser(token.Email))
		})
	secured_realm.POST("/", func(c *gin.Context) {
			ginToken, _ := c.Get("token")
			token := ginToken.(ginkeycloak.KeyCloakToken)
			crt, err := createClientForMail(token.Email)
			if(err != nil){
				c.JSON(http.StatusBadRequest, err.Error)
				return
			}
			c.JSON(http.StatusOK, crt)
	})
	secured_realm.GET("/:cn",func(c *gin.Context){
		cn := c.Param("cn")
		ginToken, _ := c.Get("token")
			token := ginToken.(ginkeycloak.KeyCloakToken)
		
		file, err := downloadCertByCN(cn, token.Email);
		if(err != nil){
			c.JSON(http.StatusBadRequest, err.Error)
			return
		}
		c.Header("Content-Type", "text/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s.ovpn\"", cn))
		c.Header("Accept-Length", fmt.Sprintf("%d", len(file)))
		c.Writer.Write([]byte(strings.Join(file, "\n")))
	})
	
	secured_realm.DELETE("/:cn",func(c *gin.Context){
		cn := c.Param("cn")
		ginToken, _ := c.Get("token")
		token := ginToken.(ginkeycloak.KeyCloakToken)
		err := revokeCertByCN(cn, token.Email)
		if(err != nil){
			c.JSON(http.StatusBadRequest, err.Error)
			return
		}
		c.JSON(http.StatusOK, token)
	})
	
}
func RunWebApi(port int, db *DB){
	router := mux.NewRouter()
	
	router.HandleFunc("/api/accessgroup", db.loadAllAccessGroups).Methods("GET")
	router.HandleFunc("/api/accessgroup", db.createAccessGroup).Methods("POST")
	router.HandleFunc("/api/accessgroup/{id}", db.deleteAccessGroup).Methods("DELETE")
	router.HandleFunc("/api/accessgroup/{id}", db.updateAccessGroup).Methods("PUT")

	router.HandleFunc("/api/user", db.loadAllUsers).Methods("GET")
	router.HandleFunc("/api/user", db.createUser).Methods("POST")
	router.HandleFunc("/api/user/{id}", db.deleteUser).Methods("DELETE")
	router.HandleFunc("/api/user/{id}", db.updateUser).Methods("PUT")

	router.HandleFunc("/api/user/{id}/access", db.loadUserAccess).Methods("GET")
	router.HandleFunc("/api/user/{id}/access", db.createUserAccess).Methods("POST")
	router.HandleFunc("/api/user/{id}/access/{group}", db.deleteUserAccess).Methods("DELETE")
	
	/*
	router.HandleFunc("/api/user/{id}/certificate", db.loadUserCertificates).Methods("GET")
	router.HandleFunc("/api/user/{id}/certificate", db.createUserCertificate).Methods("POST")
	router.HandleFunc("/api/user/{id}/certificate/{cert}", db.deleteUserCertificate).Methods("DELETE")
	router.HandleFunc("/api/user/{id}/certificate/{cert}", db.downloadUserCertificate).Methods("GET")
	*/

	router.Use(loggingMiddleware)
	
	//Static content
	/*
	staticDir := "/docker/public"
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))
	*/

	log.Println("Starting normal API operation!")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("[webapi-request] %s: Begin of %s\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		log.Printf("[webapi-request] %s: End of %s\n", r.Method, r.RequestURI)
	})
}
