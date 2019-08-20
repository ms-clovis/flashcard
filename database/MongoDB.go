package database

import (
	"context"
	"github.com/ms-clovis/flashcard/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoDB struct {
	DB RealDB
}

func (m *MongoDB) InitDB(dataSourceURI string, DBName string) {
	// Set client options
	clientOptions := options.Client().ApplyURI(dataSourceURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	m.DB.MongoDBConnection = client.Database(DBName)

}

func (m *MongoDB) GetDBUsers() ([]domain.User, bool) {
	userSlc := []domain.User{}
	collection := m.DB.MongoDBConnection.Collection("users")
	curr, err := collection.Find(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return userSlc, false
	}
	for curr.Next(context.TODO()) {
		u := domain.User{}
		err = curr.Decode(&u)
		if err != nil {
			log.Fatal(err)
			return userSlc, false
		}
		userSlc = append(userSlc, u)
	}
	return userSlc, true

}

func (m *MongoDB) GetUserFromDB(userName string) (domain.User, bool) {
	var u domain.User
	coll := m.DB.MongoDBConnection.Collection("users")
	filter := bson.D{{"username", userName}}
	err := coll.FindOne(context.TODO(), filter).Decode(&u)

	if err != nil {

		log.Println(err)
		return u, false
	}
	return u, true
}

func (m *MongoDB) AddUsersToDB(users []domain.User) bool {
	interUsers := make([]interface{}, len(users))
	ctr := 0
	for _, us := range users {

		interUsers[ctr] = us
		ctr++
	}

	coll := m.DB.MongoDBConnection.Collection("users")
	result, err := coll.InsertMany(context.TODO(), interUsers)
	if err != nil {

		log.Println(err)
		return false
	}
	return len(result.InsertedIDs) == len(users)

}

func (m *MongoDB) UpdateUserInDB(user domain.User) (domain.User, bool) {
	filter := bson.D{{"username", user.UserName}}
	coll := m.DB.MongoDBConnection.Collection("users")
	update := bson.M{
		"$set": bson.M{
			"password":  user.Password,
			"firstname": user.FirstName,
			"lastname":  user.LastName,
			"roles":     user.Roles,
		},
	}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		return domain.User{}, false
	}
	return user, result.ModifiedCount == 1

}

func (m *MongoDB) DeleteUserFromDB(user domain.User) bool {
	filter := bson.D{{
		"username", user.UserName,
	}}
	coll := m.DB.MongoDBConnection.Collection("users")
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return false
	}
	return result.DeletedCount == 1
}

func (m *MongoDB) GetDB() RealDB {
	return m.DB
}

func (m *MongoDB) SetDB(DB RealDB) {
	m.DB = DB
}

func (m *MongoDB) SetRoleOfUsersInDB(users []domain.User) bool {
	modifiedCount := 0
	for _, user := range users {
		filter := bson.D{{
			"username", user.UserName,
		}}
		update := bson.M{
			"$set": bson.M{
				"roles": user.Roles,
			},
		}
		coll := m.DB.MongoDBConnection.Collection("users")
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
			return false
		}
		modifiedCount++
	}

	return modifiedCount == len(users)
}
