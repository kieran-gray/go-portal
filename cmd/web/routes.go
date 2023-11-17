package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	dt "github.com/kieran-gray/go-portal/pkg/types"
	utils "github.com/kieran-gray/go-portal/pkg/utils"
)

var indexTemplateFunctions = template.FuncMap{
	"getHighestPriorityUrl":               utils.GetHighestPriorityUrl,
	"generateServiceCardData":             utils.GenerateServiceCardData,
	"generateServiceCardTabData":          utils.GenerateServiceCardTabData,
	"generateServiceCardFooterButtonData": utils.GenerateServiceCardFooterButtonData,
	"getFormattedTimeSince":               utils.GetFormattedTimeSince,
	"sortedByPriority":                    utils.SortedByPriority,
	"getRepoType":                         utils.GetRepoType,
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
	pipelineFile := app.getPipelineFile(app.config.PIPELINE_DATA_FILENAME)
	favourites := utils.GetFavourites(r)

	templateData := utils.GenerateIndexData(utils.GetDisplayServices(servicesFile.Services, favourites), pipelineFile)
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
	err = app.s3Client.Upload(app.config.SERVICES_FILENAME, servicesFile)
	if err != nil {
		app.serverError(w, err)
		messages = append(messages, dt.Message{Status: "failure", Message: "Failed to upload changes to S3"})
	} else {
		messages = append(messages, dt.Message{Status: "success", Message: "Successfully uploaded changes to S3"})
	}
	app.fileCache[app.config.SERVICES_FILENAME] = servicesFile
	app.render(w, r, "admin", utils.GenerateAdminData(servicesFile, messages...))
}

func (app *application) adminAddEnvironment(w http.ResponseWriter, r *http.Request) {
	serviceId := chi.URLParam(r, "id")
	serviceType := chi.URLParam(r, "serviceType")
	servicesFile, err := utils.ParseServicesFromRequest(r)
	if err != nil {
		app.serverError(w, err)
	}
	servicesFile.Services = utils.AddEnvironment(servicesFile.Services, serviceId, serviceType)
	app.render(w, r, "admin", utils.GenerateAdminData(servicesFile))
}

func (app *application) adminAddService(w http.ResponseWriter, r *http.Request) {
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

func (app *application) pipelines(w http.ResponseWriter, r *http.Request) {
	PipelineFile := app.getPipelineFile(app.config.PIPELINE_DATA_FILENAME)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(PipelineFile)
	if err != nil {
		app.errorLog.Print(err.Error())
	}
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

	router.Get("/services.json", app.services)
	router.Get("/pipelines.json", app.pipelines)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return router
}
