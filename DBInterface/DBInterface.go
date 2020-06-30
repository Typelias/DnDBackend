package dbinterface

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//User is used to parse data from DB
type User struct {
	Username string
	Password string
	UserRole string
}

//DBInterface handles connections to the MongoDB database
type DBInterface struct {
	client   *mongo.Client
	users    *mongo.Collection
	campains *mongo.Collection
}

//Init creates the DB interface
func (db *DBInterface) Init() {
	clientOptions := options.Client().ApplyURI("mongodb://typelias.se:27017").SetAuth(options.Credential{
		AuthSource: "DnDDB", Username: os.Getenv("MongoUser"), Password: os.Getenv("MongoPassword"),
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	db.client = client

	db.users = client.Database("DnDDB").Collection("users")
	db.campains = client.Database("DnDDB").Collection("campains")
}

func (db *DBInterface) findName(name string) bool {

	allUsers := db.GetAllUsers()

	for _, v := range allUsers {
		if v == name {
			return true
		}
	}

	return false
}

//DeleteUser deletes a user based on username
func (db *DBInterface) DeleteUser(name string) bool {
	result, err := db.users.DeleteOne(context.TODO(), bson.M{"username": name})
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(result)
	return true
}

//UpdateUser updates a users information
func (db *DBInterface) UpdateUser(user User, userToUpdate string) bool {

	result, err := db.users.ReplaceOne(context.TODO(), bson.M{"username": userToUpdate}, user)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(result)

	return true
}

//CheckUser checks user credentials
func (db *DBInterface) CheckUser(username, password string) (string, bool) {
	filter := bson.D{{"username", username}}
	var res User
	err := db.users.FindOne(context.TODO(), filter).Decode(&res)

	if err != nil {
		fmt.Println(err)
		return "", false
	}

	if res.Password == password {
		return res.UserRole, true
	}

	return "", false
}

//AddUser adds user to database
func (db *DBInterface) AddUser(username, password, userRole string) bool {
	newUser := User{username, password, userRole}

	if !db.findName(username) {
		insRes, err := db.users.InsertOne(context.TODO(), newUser)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Inserted a user: ", insRes.InsertedID)
		return true
	}

	return false
}

//GetAllUsers returns a string of all users
func (db *DBInterface) GetAllUsers() []string {
	var results []string
	cur, err := db.users.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		fmt.Println(err)
	}

	for cur.Next(context.TODO()) {
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, elem.Username)
	}

	return results

}
