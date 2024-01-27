package main

import (
	"flag"
	"os"

	"github.com/joho/godotenv"

	S3Client "github.com/kieran-gray/go-portal/pkg/s3Client"
	dt "github.com/kieran-gray/go-portal/pkg/types"
	utils "github.com/kieran-gray/go-portal/pkg/utils"
)

const mockFilesPath = "./data/"

type scriptConfig struct {
	SERVICES_FILENAME      string
	WORKFLOW_DATA_FILENAME string
	BUCKET_CONFIG          S3Client.BucketConfig
}

func main() {
	envFilePtr := flag.String("env-file", "/config/.env.local", "Path to env file")
	flag.Parse()

	err := godotenv.Load(*envFilePtr)
	if err != nil {
		panic(err)
	}
	config := scriptConfig{
		SERVICES_FILENAME:      utils.EnsureEnv("SERVICES_FILENAME"),
		WORKFLOW_DATA_FILENAME: utils.EnsureEnv("WORKFLOW_DATA_FILENAME"),
		BUCKET_CONFIG: S3Client.BucketConfig{
			AWS_ACCESS_KEY:              utils.EnsureEnv("AWS_ACCESS_KEY"),
			AWS_SECRET_KEY:              utils.EnsureEnv("AWS_SECRET_KEY"),
			AWS_BUCKET:                  utils.EnsureEnv("AWS_BUCKET"),
			AWS_REGION:                  utils.EnsureEnv("AWS_REGION"),
			AWS_ENDPOINT:                os.Getenv("AWS_ENDPOINT"),
			AWS_USE_PATH_STYLE_ENDPOINT: utils.ParseEnvToBool("AWS_USE_PATH_STYLE_ENDPOINT", false),
			AWS_DISABLE_SSL:             utils.ParseEnvToBool("AWS_DISABLE_SSL", false),
		},
	}

	s3Client := S3Client.S3Client(config.BUCKET_CONFIG)

	var servicesFile dt.ServicesFile
	servicesFile, err = utils.ReadFile[dt.ServicesFile](mockFilesPath+config.SERVICES_FILENAME, servicesFile)
	if err != nil {
		panic(err)
	}
	err = s3Client.Upload(config.SERVICES_FILENAME, servicesFile)
	if err != nil {
		panic(err)
	}

	var workflowFile dt.WorkflowFile
	workflowFile, err = utils.ReadFile[dt.WorkflowFile](mockFilesPath+config.WORKFLOW_DATA_FILENAME, workflowFile)
	if err != nil {
		panic(err)
	}
	s3Client.Upload(config.WORKFLOW_DATA_FILENAME, workflowFile)
	if err != nil {
		panic(err)
	}
}
