package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cassiuspaim/portimporter/domain/entities"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// #todo remove global variable
var dbClient *mongo.Client

const MongoInitdbRootUsername = "root"
const MongoInitdbRootPassword = "password"

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image for version 5.0
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			// username and password for mongodb superuser
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}

		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// run tests
	code := m.Run()

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client
	if err = dbClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestFindPorts(t *testing.T) {
	t.Parallel()
	t.Run("Given an id of a non existing Port When GetByID is invoked Then no error and not Port is expected", func(t *testing.T) {
		t.Parallel()

		portRepository := NewPortRepository(dbClient, "portsTest")

		port, err := portRepository.GetByID("idunique")
		assert.Nil(t, port, "Port must not exist at database")
		assert.NoError(t, err, "Error must not be found")
	})

	t.Run("Given a Port is created When the Port is queried it must be found", func(t *testing.T) {
		t.Parallel()

		portRepository := NewPortRepository(dbClient, "portsTest")

		err := portRepository.Create(entities.NewPort(
			"id",
			"name",
			"city",
			"country",
			[]string{"alias1", "alias2"},
			[]string{"region1", "region2"},
			[]float64{43.434343434, 35.2423434},
			"province",
			"timezone",
			[]string{"unloc1", "unloc2"},
			"code"))
		assert.NoError(t, err, "Error must not be found creating Port")

		port, err := portRepository.GetByID("id")
		assert.NotNil(t, port, "Port must exist at database")
		assert.NoError(t, err, "Error must not be found quering Port")
	})

	t.Run("Given a Port is stored When the Port is queried it must be found", func(t *testing.T) {
		t.Parallel()

		portRepository := NewPortRepository(dbClient, "portsTest")
		idPort := "idx"

		port := entities.NewPort(
			idPort,
			"name",
			"city",
			"country",
			[]string{"alias1", "alias2"},
			[]string{"region1", "region2"},
			[]float64{43.434343434, 35.2423434},
			"province",
			"timezone",
			[]string{"unloc1", "unloc2"},
			"code")
		err := portRepository.Create(port)
		assert.NoError(t, err, "Error must not be found creating Port")

		expectedCity := "Other city"
		port.City = expectedCity
		err = portRepository.Update(port, idPort)
		assert.NoError(t, err, "Error must not be found quering Port")
		portExisting, _ := portRepository.GetByID(idPort)
		assert.Equal(t, expectedCity, portExisting.City)
	})
}
