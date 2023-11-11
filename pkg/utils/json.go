package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"

	dt "github.com/kieran-gray/go-portal/pkg/types"
)

func ReadServices() (dt.ServicesFile, error) {
	var servicesFile dt.ServicesFile

	file, err := os.Open("./services.json")
	if err != nil {
		return servicesFile, err
	}
	defer file.Close()

	servicesBytes, _ := io.ReadAll(file)

	err = json.Unmarshal(servicesBytes, &servicesFile)
	if err != nil {
		return servicesFile, err
	}

	return servicesFile, nil
}

func ReadPipelineData() (dt.PipelineFile, error) {
	var pipelineFile dt.PipelineFile

	file, err := os.Open("./gitlabPipelineData.json")
	if err != nil {
		return pipelineFile, err
	}
	defer file.Close()

	pipelinesBytes, _ := io.ReadAll(file)

	err = json.Unmarshal(pipelinesBytes, &pipelineFile)
	if err != nil {
		return pipelineFile, err
	}

	return pipelineFile, nil
}

func SaveServices(servicesFile dt.ServicesFile) {
	file, err := json.MarshalIndent(servicesFile, "", " ")
	if err != nil {
		log.Print(err.Error())
		return
	}
	_ = os.WriteFile("services_output.json", file, 0644)
}
