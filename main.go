package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	dndinterface "github.com/Typelias/DnDBackend/DBInterface"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var jwtKey = []byte(os.Getenv("JWTKey"))

var users = map[string]string{
	"typelias": "pass",
	"rass":     "pass",
}

var db dndinterface.DBInterface

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

	userRole, authorized := db.CheckUser(creds.Username, creds.Password)

	if !authorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	experationTime := time.Now().Add(48 * time.Hour)

	claims := &Claims{
		Username: creds.Username,
		Type:     userRole,
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

func addUser(w http.ResponseWriter, r *http.Request) {
	var user dndinterface.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := db.AddUser(user.Username, user.Password, user.UserRole)

	if res {
		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getUserList(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(db.GetAllUsers())
}

type userDeletePost struct {
	Username string `json:"username"`
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var username userDeletePost
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	res := db.DeleteUser(username.Username)

	if res {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

type userUpdatePost struct {
	UserToUpdate string            `json:"userToUpdate"`
	User         dndinterface.User `json:"user"`
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var postData userUpdatePost
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	res := db.UpdateUser(postData.User, postData.UserToUpdate)

	if res {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func main() {
	db.Init()

	temp := db.GetAllUsers()

	fmt.Println(temp)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/signin", signIn).Methods("POST", "OPTIONS")
	router.Handle("/addUser", isAuthorized(addUser)).Methods("POST", "OPTIONS")
	router.Handle("/getUserList", isAuthorized(getUserList)).Methods("GET")
	router.Handle("/deleteUser", isAuthorized(deleteUser)).Methods("POST", "OPTIONS")
	router.Handle("/updateUser", isAuthorized(updateUser)).Methods("POST", "OPTIONS")

	headers := handlers.AllowedHeaders([]string{"accept", "authorization", "content-type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	x := handlers.ExposedHeaders([]string{"Set-Cookie"})
	cred := handlers.AllowCredentials()

	fmt.Println("Server started")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins, x, cred)(router)))

}
