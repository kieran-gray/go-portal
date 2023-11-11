package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

import dt "github.com/kieran-gray/go-portal/pkg/types"

func GenerateIndexData(displayServices dt.DisplayServices, pipelineFile dt.PipelineFile) dt.IndexData {
	return dt.IndexData{
		Services:   displayServices.Services,
		Favourites: displayServices.Favourites,
		Pipelines:  pipelineFile.Pipelines,
	}
}

func GenerateServiceCardData(service dt.Service, favourite bool, pipelines map[string]dt.Pipeline) dt.ServiceCardData {
	return dt.ServiceCardData{
		Service:       service,
		Id:            GenerateServiceId(service),
		Favourite:     favourite,
		SearchAliases: strings.Join([]string{service.Metadata.Name, service.Metadata.Aliases}, ", "),
		HasUi:         len(service.Ui.Environments) > 0,
		HasApi:        len(service.Api.Environments) > 0,
		Pipelines:     pipelines,
	}
}

func GenerateAdminCardData(service dt.Service) dt.AdminCardData {
	return dt.AdminCardData{
		Service:       service,
		Id:            GenerateServiceId(service),
		SearchAliases: strings.Join([]string{service.Metadata.Name, service.Metadata.Aliases}, ", "),
	}
}

func GenerateAdminCardTabData(service dt.Service, serviceType string) dt.AdminCardTabData {
	var serviceDetails dt.ServiceDetails
	if serviceType == "ui" {
		serviceDetails = service.Ui
	} else {
		serviceDetails = service.Api
	}
	return dt.AdminCardTabData{
		ServiceDetails: serviceDetails,
		Name:           service.Metadata.Name,
		ServiceType:    serviceType,
		Id:             GenerateServiceId(service),
	}
}

func generatePipelineData(serviceDetails dt.ServiceDetails, pipelines map[string]dt.Pipeline, serviceId string) dt.PipelineData {
	pipelineData := dt.PipelineData{
		Pipeline:    dt.Pipeline{},
		HasPipeline: false,
	}

	pipeline, ok := pipelines[serviceDetails.RepositoryUrl]
	if ok {
		pipelineData.Pipeline = pipeline
		pipelineData.HasPipeline = true
	}
	return pipelineData
}

func GenerateServiceCardTabData(service dt.Service, serviceType string, favourite bool, pipelines map[string]dt.Pipeline) dt.ServiceCardTabData {
	serviceId := GenerateServiceId(service)
	showTab := false

	var serviceDetails dt.ServiceDetails
	if serviceType == "ui" {
		serviceDetails = service.Ui
		showTab = true
	} else {
		serviceDetails = service.Api
		if len(service.Ui.Environments) == 0 {
			showTab = true
		}
	}

	return dt.ServiceCardTabData{
		ServiceDetails: serviceDetails,
		Id:             serviceId,
		Name:           service.Metadata.Name,
		ServiceType:    serviceType,
		HasLogs:        HasLogs(serviceDetails),
		Favourite:      favourite,
		Pipeline:       generatePipelineData(serviceDetails, pipelines, serviceId),
		ShowTab:        showTab,
	}
}

func GenerateServiceCardFooterButtonData(environment dt.Environment, hasLogs bool) dt.ServiceCardFooterButtonData {
	caser := cases.Caser(cases.Title(language.BritishEnglish))
	return dt.ServiceCardFooterButtonData{
		HasLogs:         hasLogs,
		EnvironmentName: caser.String(environment.Name),
		Url:             environment.Url,
		LogsUrl:         environment.LogsUrl,
	}
}
