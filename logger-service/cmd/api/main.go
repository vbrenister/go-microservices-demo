package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net/http"
	"os"
	"time"

	"github.com/vbrenister/apicommon"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://localhost:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models

	apicommon.ServerConfig
}

func main() {
	log.Printf("Starting logging services at port %s", webPort)

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := &Config{
		Models: data.New(client),
	}

	log.Panic(app.serve())
}

func (c *Config) serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: c.routes(),
	}

	return srv.ListenAndServe()
}

func connectToMongo() (*mongo.Client, error) {
	var retries int

	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	mongoURL := fmt.Sprintf("mongodb://%s", os.Getenv("MONGO_URL"))

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	for {
		c, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Println("Mongo not yet ready")
			retries++
		} else {
			log.Println("Mongo is ready")
			return c, nil

		}

		if retries > 10 {
			log.Println(err)
			return nil, err
		}

		log.Println("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
