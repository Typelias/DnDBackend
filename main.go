package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"typelias": "pass",
	"rass":     "pass",
}

//Credentials is used to parse incoming login data
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims is used to set claims when token is created
type Claims struct {
	Username string `json:"username"`
	Type     string `json:"Type"`
	jwt.StandardClaims
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	w.Header().Set("Access-Control-Allow-Credentials", "true")

	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var UserType string

	if creds.Username == "typelias" {
		UserType = "Admin"
	} else {
		UserType = "User"
	}

	experationTime := time.Now().Add(30 * time.Minute)

	claims := &Claims{
		Username: creds.Username,
		Type:     UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: experationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(tokenString)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: experationTime,
	})
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value

		tkn, err := jwt.Parse(tknStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		endpoint(w, r)
	})
}

func newWelcome(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("Welcome %s!", "user")))

}

func welcome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

func exampleGet(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("hello")
}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/signin", signIn).Methods("POST", "OPTIONS")
	router.HandleFunc("/welcome", welcome)
	router.Handle("/welcome2", isAuthorized(newWelcome))
	router.HandleFunc("/test", exampleGet).Methods("GET", "OPTIONS")
	router.Handle("/test2", isAuthorized(exampleGet)).Methods("GET")

	headers := handlers.AllowedHeaders([]string{"accept", "authorization", "content-type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	x := handlers.ExposedHeaders([]string{"Set-Cookie"})
	cred := handlers.AllowCredentials()

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins, x, cred)(router)))

}
