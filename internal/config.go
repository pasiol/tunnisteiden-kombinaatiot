package internal

import (
	"os"
	"strconv"

	"github.com/labstack/gommon/log"
	pq "github.com/pasiol/gopq"
	"github.com/pasiol/mongoutils"
)

func Teachers() pq.PrimusQuery {
	pq := pq.PrimusQuery{}
	pq.Charset = "UTF-8"
	pq.Database = "opthenk"
	pq.Sort = ""
	pq.Search = "(K37<>\"\")"
	pq.Data = "#DATA{V1};#DATA{K37};#DATA{K13}"
	pq.Footer = ""

	return pq
}

func ReadMongoSecrets() mongoutils.MongoURI {
	db, exists := os.LookupEnv("MONGO_DB")
	if !exists {
		log.Fatalf("Missing MONGO_DB environmental variable.")
	}

	user, exists := os.LookupEnv("MONGO_USER")
	if !exists {
		log.Fatalf("Missing MONGO_USER environmental variable.")
	}

	password, exists := os.LookupEnv("MONGO_PASSWORD")
	if !exists {
		log.Fatalf("Missing MONGO_PASSWORD environmental variable.")
	}

	cluster, exists := os.LookupEnv("MONGO_CLUSTER")
	if !exists {
		log.Fatalf("Missing MONGO_CLUSTER environmental variable.")
	}

	m := mongoutils.MongoURI{
		Db:       db,
		User:     user,
		Password: password,
		Cluster:  cluster,
	}

	return m
}

func ReadPrimusSecrets() (string, string, string, string) {
	host, exists := os.LookupEnv("PRIMUS_HOST")
	if !exists {
		log.Fatalf("Missing PRIMUS_HOST environmental variable.")
	}

	port, exists := os.LookupEnv("PRIMUS_PORT")
	if !exists {
		log.Fatalf("Missing PRIMUS_PORT environmental variable.")
	}

	username, exists := os.LookupEnv("PRIMUS_USER")
	if !exists {
		log.Fatalf("Missing PRIMUS_USER environmental variable.")
	}

	password, exists := os.LookupEnv("PRIMUS_PWD")
	if !exists {
		log.Fatalf("Missing PRIMUS_PWD environmental variable.")
	}

	return host, port, username, password
}

func ReadPGSecrets() (string, int, string, string, string) {
	host, exists := os.LookupEnv("PG_HOST")
	if !exists {
		log.Fatalf("Missing PG_HOST environmental variable.")
	}

	p, exists := os.LookupEnv("PG_PORT")
	if !exists {
		log.Fatalf("Missing PG_PORT environmental variable.")
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatalf("Converting port to int failed.")
	}

	user, exists := os.LookupEnv("PG_USER")
	if !exists {
		log.Fatalf("Missing PG_USER environmental variable.")
	}

	password, exists := os.LookupEnv("PG_PWD")
	if !exists {
		log.Fatalf("Missing PG_PWD environmental variable.")
	}

	db, exists := os.LookupEnv("PG_DB")
	if !exists {
		log.Fatalf("Missing PG_DB environmental variable.")
	}

	return host, port, user, password, db
}
