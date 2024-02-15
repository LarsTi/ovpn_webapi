package main

import(
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"crypto/x509"

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
	
	endpointCerts(r, config, db)
	
	log.Println("Starting normal API operation!")
	r.Run()
}
func endpointCerts(router *gin.Engine, config ginkeycloak.BuilderConfig, db *DB){
	secured_realm := router.Group("/api/certs")
	secured_realm.Use(ginkeycloak.NewAccessBuilder(config).
		RestrictButForRole("vpn-admin-user").
		Build())

		//Alle Zertifikate anzeigen für den Benutzer
	secured_realm.GET("/", func(c *gin.Context) {
		ginToken, _ := c.Get("token")
		token := ginToken.(ginkeycloak.KeyCloakToken)
		c.JSON(http.StatusOK, db.getCertsForUser(token.Email))
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
	
	router.HandleFunc("/api/user/{id}/certificate", db.loadUserCertificates).Methods("GET")
	router.HandleFunc("/api/user/{id}/certificate", db.createUserCertificate).Methods("POST")
	router.HandleFunc("/api/user/{id}/certificate/{cert}", db.deleteUserCertificate).Methods("DELETE")
	router.HandleFunc("/api/user/{id}/certificate/{cert}", db.downloadUserCertificate).Methods("GET")

	router.Use(loggingMiddleware)
	router.Use(authMiddleware)
	
	//Static content
	staticDir := "/docker/public"
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	log.Println("Starting normal API operation!")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			
			return x509.ParsePKCS1PublicKey([]byte("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs5hYk2XbSDDu90exKlZ9dB0tVcF5jkEMYFSsSG6d6EeT3bvefwqQ+NOUS112Mv4vZxFv3yyNeEWDw1tX0cVZ/uZte9MK8cfY5CI0bhtyc3mI1yI/nZVg1m0F/pUGLhIGSdHs8Cs7h3Yb/ufxgowilHwZxIDPN3qCr6yktPMECzUol+AmBf6HwV6MxArE1a58RC2+vt0dT7zCLx5tjx4QoBO13g7yn2t9EKRfHu4L3zLX1OE5z8kAVtA+s5OJDcr+x/Dw3vrjd0v7nDFRe8xj+hInNtT4uSeVOIbtrSjW4SXNWXfsplvM2xjAG/Jr3hXzBel8XcmR0Na+P2CzgEsP9QIDAQAB"))
		})
		
		if err != nil || !token.Valid {
			log.Printf("Error: %s", err.Error())
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		// Other checks, ISS, Aud, Expirey, etc ...
		// If needed, store the user principal 
		// and other relevant info the request context  
		log.Printf("token: %s", token)
		next.ServeHTTP(w, r)
	})
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
