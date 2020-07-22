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

//Character struct describes a DnD character
type Character struct {
	characterName                  string
	characterClass                 string
	level                          int
	exp                            int
	background                     string
	race                           string
	alignment                      string
	playerName                     string
	expPoints                      int
	inspiration                    bool
	proficiencyBonus               int
	savingThrows                   interface{}
	skills                         interface{}
	hp                             interface{}
	personality                    interface{}
	attacksAndSpellcasting         interface{}
	passiveInvestigation           int
	passivePerception              int
	passiveInsight                 int
	otherProficienciesAndLanguages interface{}
	equipment                      interface{}
	featuresAndTraits              interface{}
	spellcastingAbility            string
	spellSaveDC                    int
	spellAttackBonus               int
	spellList                      interface{}
	classAttributes                []string
}

//Campaign is used to handle operation on campain collection
type Campaign struct {
	Name       string
	DM         string
	Players    []string
	Characters []string
	Image      string
}

//DBInterface handles connections to the MongoDB database
type DBInterface struct {
	client     *mongo.Client
	users      *mongo.Collection
	campains   *mongo.Collection
	characters *mongo.Collection
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

// AddCharacter adds Character to the database and adds it to a campaign
func (db *DBInterface) AddCharacter(campaignName string, character interface{}) bool {
	insRes, err := db.characters.InsertOne(context.TODO(), character)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Inserted a Character: ", insRes.InsertedID)

	filter := bson.D{{"name", campaignName}}

	var camp Campaign

	db.campains.FindOne(context.TODO(), filter).Decode(&camp)

	camp.Characters = append(camp.Characters, insRes.InsertedID.(string))

	db.campains.ReplaceOne(context.TODO(), bson.M{"name": campaignName}, camp)
	return true

}

//GetCharacterByID gets a character based on an ID
func (db *DBInterface) GetCharacterByID(id string) (Character, bool) {
	filter := bson.D{{"_id", id}}
	var res Character
	err := db.users.FindOne(context.TODO(), filter).Decode(&res)

	if err != nil {
		fmt.Println(err)
		return Character{}, false
	}

	return res, true
}

//UpdateCharacter updates a character given an ID
func (db *DBInterface) UpdateCharacter(id string, ch interface{}) bool {
	filter := bson.D{{"_id", id}}
	res, err := db.characters.ReplaceOne(context.TODO(), filter, ch)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(res)
	return true
}

//RemoveCharacter removes a character based on ID
func (db *DBInterface) RemoveCharacter(id string) bool {
	filter := bson.D{{"_id", id}}
	res, err := db.characters.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(res)

	return true
}

//AddCampain adds new campains to the database
func (db *DBInterface) AddCampain(campain Campaign) bool {
	if !db.findCampain(campain.Name) {
		insRes, err := db.campains.InsertOne(context.TODO(), campain)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Inserted a campain: ", insRes.InsertedID)
		return true

	}

	return false
}

//UpdateCampaign is used to update a campaign
func (db *DBInterface) UpdateCampaign(name string, campaignToUpdate Campaign) bool {

	filter := bson.D{{"name", name}}

	var oldeVersion Campaign

	db.campains.FindOne(context.TODO(), filter).Decode(&oldeVersion)

	campaignToUpdate.Characters = oldeVersion.Characters

	result, err := db.campains.ReplaceOne(context.TODO(), bson.M{"name": name}, campaignToUpdate)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(result)

	return true
}

//RemoveCampaign removes a Campaign based on username
func (db *DBInterface) RemoveCampaign(name string) bool {
	filter := bson.D{{"name", name}}

	var oldeVersion Campaign

	db.campains.FindOne(context.TODO(), filter).Decode(&oldeVersion)

	list := oldeVersion.Characters

	for _, v := range list {
		db.RemoveCharacter(v)
	}

	result, err := db.campains.DeleteOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(result)
	return true
}

//GetUserCampaign gets specific user campaigns
func (db *DBInterface) GetUserCampaign(username string) []Campaign {
	var results []Campaign
	cur, err := db.campains.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		fmt.Println(err)
	}

	for cur.Next(context.TODO()) {
		var elem Campaign
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}
		if checkForUser(username, elem.Players) {
			results = append(results, elem)

		}
	}

	return results
}

//GetDMCampaign gets specific campaigns for a specific DM
func (db *DBInterface) GetDMCampaign(username string) []Campaign {
	var results []Campaign
	cur, err := db.campains.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		fmt.Println(err)
	}

	for cur.Next(context.TODO()) {
		var elem Campaign
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}
		if username == elem.DM {
			results = append(results, elem)

		}
	}

	return results
}

func checkForUser(username string, list []string) bool {

	for _, v := range list {
		if v == username {
			return true
		}

	}

	return false
}

//GetAllCampains gets alla campains
func (db *DBInterface) GetAllCampains() []Campaign {
	var results []Campaign
	cur, err := db.campains.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		fmt.Println(err)
	}

	for cur.Next(context.TODO()) {
		var elem Campaign
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, elem)
	}

	return results
}

func (db *DBInterface) findCampain(name string) bool {
	allCampains := db.GetAllCampains()

	for _, v := range allCampains {
		if v.Name == name {
			return true
		}
	}

	return false

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
