package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	dt "github.com/kieran-gray/go-portal/pkg/types"
	"github.com/kieran-gray/go-portal/pkg/utils"
)

var indexTemplateFunctions = template.FuncMap{
	"getHighestPriorityUrl":               utils.GetHighestPriorityUrl,
	"generateServiceCardData":             utils.GenerateServiceCardData,
	"generateServiceCardTabData":          utils.GenerateServiceCardTabData,
	"generateServiceCardFooterButtonData": utils.GenerateServiceCardFooterButtonData,
	"getFormattedTimeSince":               utils.GetFormattedTimeSince,
	"sortedByPriority":                    utils.SortedByPriority,
	"getWorkflowStatus":                   utils.GetWorkflowStatus,
}

var adminTemplateFunctions = template.FuncMap{
	"generateAdminCardData":    utils.GenerateAdminCardData,
	"toMap":                    utils.ToMap,
	"generateAdminCardTabData": utils.GenerateAdminCardTabData,
}

var templateToFuncMap = map[string]template.FuncMap{
	"index": indexTemplateFunctions,
	"admin": adminTemplateFunctions,
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	servicesFile := app.getServicesFile(app.config.SERVICES_FILENAME)
	workflowFile := app.getWorkflowFile(app.config.WORKFLOW_DATA_FILENAME)
	favourites := utils.GetFavourites(r)

	templateData := utils.GenerateIndexData(
		utils.GetDisplayServices(servicesFile.Services, favourites),
		workflowFile,
	)
	app.render(w, r, "index", templateData)
}

func (app *application) adminGet(w http.ResponseWriter, r *http.Request) {
	servicesFile := app.getServicesFile(app.config.SERVICES_FILENAME)
	app.render(w, r, "admin", utils.GenerateAdminData(servicesFile))
}

func (app *application) adminPost(w http.ResponseWriter, r *http.Request) {
	messages := []dt.Message{}
	servicesFile, err := utils.ParseServicesFromRequest(r)
	if err != nil {
		app.serverError(w, err)
	}
	servicesFile.Services = utils.SortServicesByName(servicesFile.Services)
	err = app.s3Client.Upload(app.config.SERVICES_FILENAME, servicesFile)
	if err != nil {
		app.serverError(w, err)
		messages = append(
			messages,
			dt.Message{
				Status:  "failure",
				Message: "Failed to upload changes to S3",
			},
		)
	} else {
		messages = append(
			messages,
			dt.Message{
				Status:  "success",
				Message: "Successfully uploaded changes to S3",
			},
		)
	}
	app.fileCache[app.config.SERVICES_FILENAME] = servicesFile
	app.render(
		w, r, "admin", utils.GenerateAdminData(servicesFile, messages...),
	)
}

func (app *application) adminAddEnvironment(
	w http.ResponseWriter, r *http.Request,
) {
	serviceId := chi.URLParam(r, "id")
	serviceType := chi.URLParam(r, "serviceType")
	servicesFile, err := utils.ParseServicesFromRequest(r)
	if err != nil {
		app.serverError(w, err)
	}
	servicesFile.Services = utils.AddEnvironment(
		servicesFile.Services, serviceId, serviceType,
	)
	app.render(w, r, "admin", utils.GenerateAdminData(servicesFile))
}

func (app *application) adminAddService(
	w http.ResponseWriter, r *http.Request,
) {
	servicesFile, err := utils.ParseServicesFromRequest(r)
	if err != nil {
		app.serverError(w, err)
	}
	servicesFile.Services = utils.AddService(servicesFile.Services)
	app.render(w, r, "admin", utils.GenerateAdminData(servicesFile))
}

func (app *application) services(w http.ResponseWriter, r *http.Request) {
	servicesFile := app.getServicesFile(app.config.SERVICES_FILENAME)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(servicesFile)
	if err != nil {
		app.errorLog.Print(err.Error())
	}
}

func (app *application) workflows(w http.ResponseWriter, r *http.Request) {
	workflowFile := app.getWorkflowFile(app.config.WORKFLOW_DATA_FILENAME)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(workflowFile)
	if err != nil {
		app.errorLog.Print(err.Error())
	}
}

func (app *application) workflowWebhook(
	w http.ResponseWriter, r *http.Request,
) {
	webhookWorkflow, err := utils.ParseWebhookWorkflowFromRequest(r)
	if err != nil {
		app.errorLog.Print("Failed to read webhook request body")
		app.serverError(w, err)
	}
	repo := webhookWorkflow.Repository
	if webhookWorkflow.WorkflowRun.HeadBranch != repo.DefaultBranch {
		// ignore workflows from non-main branch
		app.okResponse(w)
		return
	}

	workflowFile := app.getWorkflowFile(app.config.WORKFLOW_DATA_FILENAME)
	if workflowFile.Workflows == nil {
		workflowFile.Workflows = make(map[string]dt.Workflow)
	}
	workflowFile.Workflows[repo.Url] = utils.WebhookWorkflowToWorkflow(
		webhookWorkflow,
	)
	err = app.s3Client.Upload(app.config.WORKFLOW_DATA_FILENAME, workflowFile)
	if err != nil {
		app.serverError(w, err)
	}
	app.fileCache[app.config.WORKFLOW_DATA_FILENAME] = workflowFile
	app.createdResponse(w)
}

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", app.index)
	router.Route("/admin", func(router chi.Router) {
		router.Get("/", app.adminGet)
		router.Post("/", app.adminPost)
		router.Post("/addEnv/{id}&{serviceType}", app.adminAddEnvironment)
		router.Post("/addService", app.adminAddService)
	})
	router.Route("/webhook", func(router chi.Router) {
		router.Post("/workflow", app.workflowWebhook)
	})

	router.Get("/services.json", app.services)
	router.Get("/workflows.json", app.workflows)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return router
}
