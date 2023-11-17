package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	dt "github.com/kieran-gray/go-portal/pkg/types"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	template, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buffer := new(bytes.Buffer)

	err := template.Execute(buffer, data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	_, err = buffer.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) getServicesFile(filename string) dt.ServicesFile {
	var servicesFile dt.ServicesFile
	inter, ok := app.fileCache[filename]
	if ok {
		servicesFile, ok := inter.(dt.ServicesFile)
		if !ok {
			app.errorLog.Print("Failed converting cached services file to struct")
			return servicesFile
		}
		return servicesFile
	}

	bytes, err := app.s3Client.Download(filename)
	if err != nil {
		app.errorLog.Printf("Failed to fetch file %s from S3. %v", filename, err)
		return servicesFile
	}
	err = json.Unmarshal(bytes, &servicesFile)
	if err != nil {
		app.errorLog.Print(err.Error())
		return servicesFile
	}
	app.fileCache[filename] = servicesFile
	return servicesFile
}

func (app *application) getPipelineFile(filename string) dt.PipelineFile {
	var pipelineFile dt.PipelineFile
	inter, ok := app.fileCache[filename]
	if ok {
		pipelineFile, ok := inter.(dt.PipelineFile)
		if !ok {
			app.errorLog.Print("Failed converting cached pipeline file to struct")
			return pipelineFile
		}
		return pipelineFile
	}
	bytes, err := app.s3Client.Download(filename)
	if err != nil {
		app.errorLog.Printf("Failed to fetch file %s from S3. %v", filename, err)
		return pipelineFile
	}

	err = json.Unmarshal(bytes, &pipelineFile)
	if err != nil {
		app.errorLog.Print(err.Error())
		return pipelineFile
	}
	app.fileCache[filename] = pipelineFile
	return pipelineFile
}
