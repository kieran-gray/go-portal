package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

import dt "github.com/kieran-gray/go-portal/pkg/types"

func GenerateIndexData(
	displayServices dt.DisplayServices, workflowFile dt.WorkflowFile,
) dt.IndexData {
	return dt.IndexData{
		Services:   displayServices.Services,
		Favourites: displayServices.Favourites,
		Workflows:  workflowFile.Workflows,
	}
}

func GenerateAdminData(
	servicesFile dt.ServicesFile, messages ...dt.Message,
) dt.AdminData {
	return dt.AdminData{
		ServicesFile: servicesFile,
		Messages:     messages,
	}
}

func GenerateServiceCardData(
	service dt.Service, favourite bool, workflows map[string]dt.Workflow,
) dt.ServiceCardData {
	return dt.ServiceCardData{
		Service:       service,
		Id:            GenerateServiceId(service),
		Favourite:     favourite,
		SearchAliases: strings.Join([]string{service.Metadata.Name, service.Metadata.Aliases}, ", "),
		HasUi:         len(service.Ui.Environments) > 0,
		HasApi:        len(service.Api.Environments) > 0,
		Workflows:     workflows,
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

func generateWorkflowData(
	serviceDetails dt.ServiceDetails, workflows map[string]dt.Workflow, serviceId string,
) dt.WorkflowData {
	workflowData := dt.WorkflowData{
		Workflow:    dt.Workflow{},
		HasWorkflow: false,
	}

	workflow, ok := workflows[serviceDetails.RepositoryUrl]
	if ok {
		workflowData.Workflow = workflow
		workflowData.HasWorkflow = true
	}
	return workflowData
}

func GenerateServiceCardTabData(serviceCardData dt.ServiceCardData, serviceType string) dt.ServiceCardTabData {
	service := serviceCardData.Service
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
		Favourite:      serviceCardData.Favourite,
		Workflow:       generateWorkflowData(serviceDetails, serviceCardData.Workflows, serviceId),
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
