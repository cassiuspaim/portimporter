package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cassiuspaim/portimporter/domain/entities"
	"github.com/cassiuspaim/portimporter/domain/services"
	"github.com/cassiuspaim/portimporter/infrastructure/jsonstream"
	"github.com/cassiuspaim/portimporter/infrastructure/repositories/mongodb"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Handle the signals to handle graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	clientDB := connectToDatabase()

	runApp(ctx, clientDB, &waitGroup, stop)

	waitGroup.Wait()
	closeApp(clientDB)
}

func runApp(ctx context.Context, dbConnect *mongo.Client, waitGroup *sync.WaitGroup, stop context.CancelFunc) {
	portRepository := mongodb.NewPortRepository(dbConnect, os.Getenv("DB_NAME"))
	portService := services.NewPortService(portRepository)
	stream := jsonstream.NewPortStream()

	go func() {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Stopping Port import. Stopping message: %v\n", ctx.Err())

				return

			default:
				for entry := range stream.Watch() {
					if entry.Error != nil {
						log.Println(entry.Error)

						continue
					}

					port := entities.Port{
						ID:          entry.Key,
						Name:        entry.Data.Name,
						Coordinates: entry.Data.Coordinates,
						City:        entry.Data.City,
						Province:    entry.Data.Province,
						Country:     entry.Data.Country,
						Alias:       entry.Data.Alias,
						Regions:     entry.Data.Regions,
						Unlocs:      entry.Data.Unlocs,
						Timezone:    entry.Data.Timezone,
						Code:        entry.Data.Code,
					}

					err := portService.Upsert(port)
					if err != nil {
						log.Printf("Error upserting the Port %s. Error: %s", port.ID, err)
					}
				}

				stop()
			}
		}
	}()

	log.Println("Openning port file.")

	fileName := os.Getenv("PORT_JSON_PATH")
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("Error opening file %s", fileName)
	}
	defer file.Close()

	stream.Start(file)
}

func connectToDatabase() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	} else {
		log.Println("Env file found.")
	}

	dbAuthName := os.Getenv("DB_AUTHENTICATION_NAME")
	dbUserName := os.Getenv("DB_USER_NAME")
	dbPassword := os.Getenv("DB_USER_PASSWORD")
	dbURI := os.Getenv("DB_CONNECTION_URI")

	log.Printf("Connecting DB through auth name %s\n", dbAuthName)
	log.Printf("Connecting DB using user name %s\n", dbUserName)
	log.Printf("Connecting DB using password %s\n", dbPassword)
	log.Printf("Connecting DB to URI %s\n", dbURI)

	// Set credential
	credential := options.Credential{
		AuthSource: dbAuthName,
		Username:   dbUserName,
		Password:   dbPassword,
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(dbURI).SetAuth(credential)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database!")

	return client
}

func closeApp(clientDB *mongo.Client) {
	log.Println("Closing app.")
	closeConnection(clientDB)
}

func closeConnection(clientDB *mongo.Client) {
	log.Println("Start closing database connection.")

	err := clientDB.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection closed.")
}
