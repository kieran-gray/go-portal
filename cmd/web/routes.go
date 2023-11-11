package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"github.com/go-chi/chi/v5"
	utils "github.com/kieran-gray/go-portal/pkg/utils"
)

var indexTemplateFunctions = template.FuncMap{
	"getHighestPriorityUrl": utils.GetHighestPriorityUrl,
	"generateServiceCardData": utils.GenerateServiceCardData,
	"generateServiceCardTabData": utils.GenerateServiceCardTabData,
	"generateServiceCardFooterButtonData": utils.GenerateServiceCardFooterButtonData,
	"getFormattedTimeSince": utils.GetFormattedTimeSince,
	"sortedByPriority": utils.SortedByPriority,
}

var adminTemplateFunctions = template.FuncMap{
	"generateAdminCardData": utils.GenerateAdminCardData,
	"toMap": utils.ToMap,
	"generateAdminCardTabData": utils.GenerateAdminCardTabData,
}

var templateToFuncMap = map[string]template.FuncMap{
	"index": indexTemplateFunctions,
	"admin": adminTemplateFunctions,
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	servicesFile := app.getServicesFile("services.json")
	pipelineFile := app.getPipelineFile("gitlabPipelineData.json")
	favourites := utils.GetFavourites(r)

	templateData := utils.GenerateIndexData(utils.GetDisplayServices(servicesFile.Services, favourites), pipelineFile)
	app.render(w, r, "index", templateData)
}

func (app *application) adminGet(w http.ResponseWriter, r *http.Request) {
	servicesFile := app.getServicesFile("services.json")
	app.render(w, r, "admin", servicesFile)
}

func (app *application) adminPost(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Print("POST /admin: processing request")
	servicesFile, err := utils.ParseServicesFromRequest(r)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Print("POST /admin: parsed request body")
	err = app.s3Client.Upload("services.json", servicesFile)
	if err != nil {
		app.serverError(w, err)
	}

	app.fileCache["services.json"] = servicesFile
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (app *application) adminAddEnvironment(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Print("POST /admin/addEnv: processing request")
	serviceId := chi.URLParam(r, "id")
	serviceType := chi.URLParam(r, "serviceType")
	servicesFile, err := utils.ParseServicesFromRequest(r)
	if err != nil {
		app.serverError(w, err)
	}
	servicesFile.Services = utils.AddEnvironment(servicesFile.Services, serviceId, serviceType)
	app.render(w, r, "admin", servicesFile)
}

func (app *application) adminAddService(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Print("POST /admin/addService: processing request")
	servicesFile, err := utils.ParseServicesFromRequest(r)
	if err != nil {
		app.serverError(w, err)
	}
	app.infoLog.Print("POST /admin/addService: parsed request body")	
	servicesFile.Services = utils.AddService(servicesFile.Services)
	app.render(w, r, "admin", servicesFile)
}

func (app *application) services(w http.ResponseWriter, r *http.Request) {
	servicesFile := app.getServicesFile("services.json")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servicesFile)
}

func (app *application) pipelines(w http.ResponseWriter, r *http.Request) {
	PipelineFile := app.getPipelineFile("gitlabPipelineData.json")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PipelineFile)
}

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()
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