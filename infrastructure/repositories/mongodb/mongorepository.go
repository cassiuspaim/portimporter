package mongodb

import (
	"context"
	"errors"
	"log"

	"github.com/cassiuspaim/portimporter/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PortDB is used by implementation for Mongo of PortRepository
type PortDB struct {
	Key         string    `bson:"key"`
	Name        string    `bson:"name"`
	City        string    `bson:"city"`
	Country     string    `bson:"country"`
	Alias       []string  `bson:"alias"`
	Regions     []string  `bson:"regions"`
	Coordinates []float64 `bson:"coordinates"`
	Province    string    `bson:"province"`
	Timezone    string    `bson:"timezone"`
	Unlocs      []string  `bson:"unlocs"`
	Code        string    `bson:"code"`
}

// Retrieves a PortDB based on entities.Port passed by parameter.
func (p PortDB) From(port entities.Port) PortDB {
	return PortDB{
		Key:         port.ID,
		Name:        port.Name,
		City:        port.City,
		Country:     port.Country,
		Alias:       port.Alias,
		Regions:     port.Regions,
		Coordinates: port.Coordinates,
		Province:    port.Province,
		Timezone:    port.Timezone,
		Unlocs:      port.Unlocs,
		Code:        port.Code,
	}
}

// Retrieves an entities.Port based on the PortDB.
func (p PortDB) To() entities.Port {
	return entities.Port{
		ID:          p.Key,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Coordinates: p.Coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		Code:        p.Code,
	}
}

type PortRepository struct {
	client       *mongo.Client
	databaseName string
}

func NewPortRepository(client *mongo.Client, databaseName string) PortRepository {
	log.Printf("Database is %s", databaseName)

	return PortRepository{
		client:       client,
		databaseName: databaseName,
	}
}

func (p PortRepository) GetByID(id string) (*entities.Port, error) {
	portsCollection := p.client.Database(p.databaseName).Collection("ports")

	var portDB PortDB

	filter := bson.D{{Key: "key", Value: id}}
	result := portsCollection.FindOne(context.TODO(), filter)
	err := result.Decode(&portDB)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	port := portDB.To()

	return &port, nil
}

func (p PortRepository) Create(port entities.Port) error {
	portsCollection := p.client.Database(p.databaseName).Collection("ports")

	var portDB PortDB

	_, err := portsCollection.InsertOne(context.TODO(), portDB.From(port))

	return err
}

func (p PortRepository) Update(port entities.Port, id string) error {
	portsCollection := p.client.Database(p.databaseName).Collection("ports")

	var portDB PortDB

	_, err := portsCollection.ReplaceOne(context.TODO(), bson.M{"key": id}, portDB.From(port))

	return err
}
