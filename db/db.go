package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

const MongoUri = "mongodb://localhost:27017"

var Client *mongo.Client
var DB *mongo.Database
var SiteDirectory *mongo.Collection

func Init() {
	var err error
	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoUri))
	if err != nil {
		panic(err)
	} // Can't use util.Check here bcs it would create an import cycle
	DB = Client.Database("ugle")
	SiteDirectory = DB.Collection("sitedirectory")
}

type Site struct {
	Domain      string `json:domain`
	IP          string `json:ip`
	DiscordID   string `json:discordID`
	Title       string `json:title`
	Description string `json:description`
	SpecifiedDescription string `json:specifiedDescription`
	SpecifiedTags string `json:specifiedTags`
}
