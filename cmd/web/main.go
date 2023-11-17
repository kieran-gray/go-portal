package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	S3Client "github.com/kieran-gray/go-portal/pkg/s3Client"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	config        config
	templateCache map[string]*template.Template
	fileCache     map[string]interface{}
	s3Client      S3Client.S3
}

type config struct {
	HOST                   string
	PORT                   string
	SERVICES_FILENAME      string
	PIPELINE_DATA_FILENAME string
	WORKFLOW_DATA_FILENAME string
	BUCKET_CONFIG          S3Client.BucketConfig
}

const templateRootDir string = "./ui/html/"

func ensureEnv(key string) string {
	value, present := os.LookupEnv(key)
	if !present {
		panic(fmt.Sprintf("%s env variable not set", key))
	}
	return value
}

func parseEnvToBool(key string, defaultValue bool) bool {
	value, present := os.LookupEnv(key)
	if !present {
		return defaultValue
	}
	converted, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return converted
}

func getConfig() config {
	return config{
		HOST:                   ensureEnv("HOST"),
		PORT:                   ensureEnv("PORT"),
		SERVICES_FILENAME:      ensureEnv("SERVICES_FILENAME"),
		PIPELINE_DATA_FILENAME: ensureEnv("PIPELINE_DATA_FILENAME"),
		WORKFLOW_DATA_FILENAME: ensureEnv("WORKFLOW_DATA_FILENAME"),
		BUCKET_CONFIG: S3Client.BucketConfig{
			AWS_ACCESS_KEY:              ensureEnv("AWS_ACCESS_KEY"),
			AWS_SECRET_KEY:              ensureEnv("AWS_SECRET_KEY"),
			AWS_BUCKET:                  ensureEnv("AWS_BUCKET"),
			AWS_REGION:                  ensureEnv("AWS_REGION"),
			AWS_ENDPOINT:                os.Getenv("AWS_ENDPOINT"),
			AWS_USE_PATH_STYLE_ENDPOINT: parseEnvToBool("AWS_USE_PATH_STYLE_ENDPOINT", false),
			AWS_DISABLE_SSL:             parseEnvToBool("AWS_DISABLE_SSL", false),
		},
	}
}

func main() {
	config := getConfig()
	addr := flag.String("addr", fmt.Sprintf("%s:%s", config.HOST, config.PORT), "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	templateCache, err := templateCache(templateRootDir)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		config:        config,
		templateCache: templateCache,
		fileCache:     map[string]interface{}{},
		s3Client:      S3Client.S3Client(config.BUCKET_CONFIG),
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("---- Starting Server on %s ----", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
