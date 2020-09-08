package dbinterface

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

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

//Stats is a subclass of character
type Stats struct {
	Strength             int `json:"strength"`
	StrengthModifier     int `json:"strengthModifier"`
	Dexterity            int `json:"dexterity"`
	DexterityModifier    int `json:"dexterityModifier"`
	Constitution         int `json:"constitution"`
	ConstitutionModifier int `json:"constitutionModifier"`
	Intelligence         int `json:"intelligence"`
	IntelligenceModifier int `json:"intelligenceModifier"`
	Wisdom               int `json:"wisdom"`
	WisdomModifier       int `json:"wisdomModifier"`
	Charisma             int `json:"charisma"`
	CharismaModifier     int `json:"charismaModifier"`
}

//SavingThrows is a subclass of character
type SavingThrows struct {
	Strength     bool `json:"strength"`
	Dexterity    bool `json:"dexterity"`
	Constitution bool `json:"constitution"`
	Intelligence bool `json:"intelligence"`
	Wisdom       bool `json:"wisdom"`
	Charisma     bool `json:"charisma"`
}

//Skills is a subclass of character
type Skills struct {
	Acrobatics          bool `json:"acrobatics"`
	AcrobaticsBonus     int  `json:"acrobaticsBonus"`
	AnimalHandling      bool `json:"animalHandling"`
	AnimalHandlingBonus int  `json:"animalHandlingBonus"`
	Arcana              bool `json:"arcana"`
	ArcanaBonus         int  `json:"arcanaBonus"`
	Athletics           bool `json:"athletics"`
	AthleticsBonus      int  `json:"athleticsBonus"`
	Deception           bool `json:"deception"`
	DeceptionBonus      int  `json:"deceptionBonus"`
	History             bool `json:"history"`
	HistoryBonus        int  `json:"historyBonus"`
	Insight             bool `json:"insight"`
	InsightBonus        int  `json:"insightBonus"`
	Intimidation        bool `json:"intimidation"`
	IntimidationBonus   int  `json:"intimidationBonus"`
	Investigation       bool `json:"investigation"`
	InvestigationBonus  int  `json:"investigationBonus"`
	Medicine            bool `json:"medicine"`
	MedicineBonus       int  `json:"medicineBonus"`
	Nature              bool `json:"nature"`
	NatureBonus         int  `json:"natureBonus"`
	Perception          bool `json:"perception"`
	PerceptionBonus     int  `json:"perceptionBonus"`
	Performance         bool `json:"performance"`
	PerformanceBonus    int  `json:"performanceBonus"`
	Persuasion          bool `json:"persuasion"`
	PersuasionBonus     int  `json:"persuasionBonus"`
	Religion            bool `json:"religion"`
	ReligionBonus       int  `json:"religionBonus"`
	SlightOfHand        bool `json:"slightOfHand"`
	SlightOfHandBonus   int  `json:"slightOfHandBonus"`
	Stealth             bool `json:"stealth"`
	StealthBonus        int  `json:"stealthBonus"`
	Survival            bool `json:"survival"`
	SurvivalBonus       int  `json:"survivalBonus"`
}

//HP is a subclass of character
type HP struct {
	ArmorClass      int    `json:"armorClass"`
	Initiative      int    `json:"initiative"`
	Speed           int    `json:"speed"`
	MaxHP           int    `json:"maxHP"`
	CurrHP          int    `json:"currHP"`
	TempHP          int    `json:"tempHP"`
	HitDice         string `json:"hitDice"`
	NumberOfHutDice int    `json:"numberOfHutDice"`
}

//Personality is a subclass of character
type Personality struct {
	PersonalityTraits string `json:"personalityTraits"`
	Ideals            string `json:"ideals"`
	Bonds             string `json:"bonds"`
	Flaws             string `json:"flaws"`
	Backstory         string `json:"backstory"`
}

//Weapon is a subclass of character
type Weapon struct {
	Name        string `json:"name"`
	Damage      string `json:"damage"`
	AtkBonus    string `json:"atkBonus"`
	DamageType  string `json:"damageType"`
	Description string `json:"description"`
	Condition   string `json:"condition"`
	Amount      string `json:"amount"`
}

//AttacksAndSpellcasting is a subclass of character
type AttacksAndSpellcasting struct {
	Weapons []Weapon `json:"weapons"`
}

//CategoryItem is a subclass of character
type CategoryItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	ShowAmount  bool   `json:"showAmount"`
}

//Category is a subclass of character
type Category struct {
	Name  string         `json:"name"`
	Items []CategoryItem `json:"items"`
}

//OtherProficienciesAndLanguages is a subclass of character
type OtherProficienciesAndLanguages struct {
	Categories []Category `json:"categories"`
}

//Currency is a subclass of character
type Currency struct {
	Cp int `json:"cp"`
	Sp int `json:"sp"`
	Ep int `json:"ep"`
	Gp int `json:"gp"`
	Pp int `json:"PP"`
}

//Equipment is a subclass of character
type Equipment struct {
	EquipmentList []CategoryItem `json:"equipmentList"`
	Currency      Currency       `json:"currency"`
}

//FeaturesAndTraits is a subclass of character
type FeaturesAndTraits struct {
	Categories []Category `json:"categories"`
}

//Spell is a subclass of character
type Spell struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Dice          string `json:"dice"`
	DamageType    string `json:"damageType"`
	SpellRange    string `json:"spellRange"`
	Component     string `json:"component"`
	Duration      string `json:"duration"`
	CastingTime   string `json:"castingTime"`
	Concentration bool   `json:"concentration"`
	Conditions    string `json:"conditions"`
	Level         string `json:"level"`
}

//SpellList is a subclass of character
type SpellList struct {
	SpellList      []Spell `json:"spellList"`
	Lvl1SpellSlots int     `json:"lvl1SpellSlots"`
	Lvl2SpellSlots int     `json:"lvl2SpellSlots"`
	Lvl3SpellSlots int     `json:"lvl3SpellSlots"`
	Lvl4SpellSlots int     `json:"lvl4SpellSlots"`
	Lvl5SpellSlots int     `json:"lvl5SpellSlots"`
	Lvl6SpellSlots int     `json:"lvl6SpellSlots"`
	Lvl7SpellSlots int     `json:"lvl7SpellSlots"`
	Lvl8SpellSlots int     `json:"lvl8SpellSlots"`
	Lvl9SpellSlots int     `json:"lvl9SpellSlots"`

	Lvl1SpellSlutsUsed int `json:"lvl1SpellSlutsUsed"`
	Lvl2SpellSlutsUsed int `json:"lvl2SpellSlutsUsed"`
	Lvl3SpellSlutsUsed int `json:"lvl3SpellSlutsUsed"`
	Lvl4SpellSlutsUsed int `json:"lvl4SpellSlutsUsed"`
	Lvl5SpellSlutsUsed int `json:"lvl5SpellSlutsUsed"`
	Lvl6SpellSlutsUsed int `json:"lvl6SpellSlutsUsed"`
	Lvl7SpellSlutsUsed int `json:"lvl7SpellSlutsUsed"`
	Lvl8SpellSlutsUsed int `json:"lvl8SpellSlutsUsed"`
	Lvl9SpellSlutsUsed int `json:"lvl9SpellSlutsUsed"`
}

//Character struct describes a DnD character
type Character struct {
	CharacterName                  string                         `json:"characterName"`
	CharacterClass                 string                         `json:"characterClass"`
	Level                          int                            `json:"level"`
	Exp                            int                            `json:"exp"`
	Background                     string                         `json:"background"`
	Race                           string                         `json:"race"`
	Alignment                      string                         `json:"alignment"`
	PlayerName                     string                         `json:"playerName"`
	ExpPoints                      int                            `json:"expPoints"`
	Stats                          Stats                          `json:"stats"`
	Inspiration                    bool                           `json:"inspiration"`
	ProficiencyBonus               int                            `json:"proficiencyBonus"`
	SavingThrows                   SavingThrows                   `json:"savingThrows"`
	Skills                         Skills                         `json:"skills"`
	Hp                             HP                             `json:"hp"`
	Personality                    Personality                    `json:"personality"`
	AttacksAndSpellcasting         AttacksAndSpellcasting         `json:"attacksAndSpellcasting"`
	PassiveInvestigation           int                            `json:"passiveInvestigation"`
	PassivePerception              int                            `json:"passivePerception"`
	PassiveInsight                 int                            `json:"passiveInsight"`
	OtherProficienciesAndLanguages OtherProficienciesAndLanguages `json:"otherProficienciesAndLanguages"`
	Equipment                      Equipment                      `json:"equipment"`
	FeaturesAndTraits              FeaturesAndTraits              `json:"featuresAndTraits"`
	SpellcastingAbility            string                         `json:"spellcastingAbility"`
	SpellSaveDC                    int                            `json:"spellSaveDC"`
	SpellAttackBonus               int                            `json:"spellAttackBonus"`
	SpellList                      SpellList                      `json:"spellList"`
	ClassAttributes                []string                       `json:"classAttributes"`
	DMComments                     string                         `json:"DMComments"`
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
	db.characters = client.Database("DnDDB").Collection("characters")
}

// AddCharacter adds Character to the database and adds it to a campaign
func (db *DBInterface) AddCharacter(campaignName string, character Character) bool {
	insRes, err := db.characters.InsertOne(context.TODO(), character)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Inserted a Character: ", insRes.InsertedID)

	filter := bson.D{{"name", campaignName}}

	var camp Campaign

	db.campains.FindOne(context.TODO(), filter).Decode(&camp)

	id := strings.Split(fmt.Sprintf("%v", insRes.InsertedID), "\"")[1]

	camp.Characters = append(camp.Characters, id)

	db.campains.ReplaceOne(context.TODO(), bson.M{"name": campaignName}, camp)
	return true

}

//GetCharacterByID gets a character based on an ID
func (db *DBInterface) GetCharacterByID(id string) (Character, bool) {

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	var res Character
	err := db.characters.FindOne(context.TODO(), filter).Decode(&res)

	if err != nil {
		fmt.Println(err)
		return Character{}, false
	}

	return res, true
}

//MultiCharacterGetReturn is a helper struct for returning mulitple characters
type MultiCharacterGetReturn struct {
	ID        string    `json:"id"`
	Character Character `json:"character"`
}

//GetMultiCharacter gets alla the character by string array of character id:s
func (db *DBInterface) GetMultiCharacter(ids []string) []MultiCharacterGetReturn {
	var ret []MultiCharacterGetReturn
	for _, v := range ids {
		ch, found := db.GetCharacterByID(v)
		if found {
			ret = append(ret, MultiCharacterGetReturn{ID: v, Character: ch})

		}
	}

	return ret
}

//UpdateCharacter updates a character given an ID
func (db *DBInterface) UpdateCharacter(id string, ch Character) bool {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
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
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
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

//GetCampaignByName gets a campaign based on its name
func (db *DBInterface) GetCampaignByName(name string) Campaign {
	filter := bson.D{{"name", name}}
	var camp Campaign
	db.campains.FindOne(context.TODO(), filter).Decode(&camp)

	return camp

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
