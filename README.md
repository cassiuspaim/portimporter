#Port Importer

The Port Importer application has the goal to import the json file to the mongo database indicated by environment variable **PORT_JSON_PATH**.

## Tech Stack
- golang
- mongodb

## Prerequisites
To run the application locally by two ways:
1. Run the application from the source folder with access to a mongo database instance.
2. Run the entire application using containers.

### Run the application from the source folder
To run the application from the source folder you will need golang installed and mongo database available or installed at your machine.

If you do not have a mongo database available or installed you can run the command bellow to create a container with mongo database. To do this you must have docker installed at your machine, you can install docker from https://docs.docker.com/get-docker/. Only take a look at the env file to match the database URI and application DB client.
```
docker-compose -f docker-compose-mongo.yml up  
```

The docker-compose-mongo.yml file is using the environment variables defined at .env file. This same file will be used by Port Importer application.

The golang can be installed from https://go.dev/doc/install.

### Run the entire application using containers
For this option you must have docker installed, to install it go to https://docs.docker.com/get-docker/.
To run the application using containers you can run the command.
```
docker-compose up
```

## Development environment
To develop the code it is configured at the repository the settings for golangci-lint. To install it you can follow the instructions at https://golangci-lint.run/usage/install/#local-installation.
If you want to change the golangci-lint configurations you look at https://golangci-lint.run/usage/configuration/.



