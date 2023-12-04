package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	S3Client "github.com/kieran-gray/go-portal/pkg/s3Client"
	dt "github.com/kieran-gray/go-portal/pkg/types"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func (app *application) okResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (app *application) createdResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func (app *application) render(
	w http.ResponseWriter, r *http.Request, name string, data interface{},
) {
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

func getFileFromS3[T any](
	s3Client S3Client.S3, filename string, t T,
) (T, error) {
	bytes, err := s3Client.Download(filename)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func getFileFromCache[T any](
	fileCache map[string]interface{}, filename string, t T,
) (T, error) {
	inter, ok := fileCache[filename]
	if !ok {
		return t, fmt.Errorf("Filename %s not found in filecache", filename)
	}
	t, ok = inter.(T)
	if !ok {
		return t, fmt.Errorf("Could not parse filename %s to type", filename)
	}
	return t, nil
}

func getFile[T any](app *application, filename string, t T) T {
	t, err := getFileFromCache[T](app.fileCache, filename, t)
	if err != nil {
		app.errorLog.Print(err.Error())
	} else {
		return t
	}
	t, err = getFileFromS3[T](app.s3Client, filename, t)
	if err != nil {
		app.errorLog.Print(err.Error())
	} else {
		app.fileCache[filename] = t
	}
	return t
}

func (app *application) getServicesFile(filename string) dt.ServicesFile {
	var servicesFile dt.ServicesFile
	return getFile[dt.ServicesFile](app, filename, servicesFile)
}

func (app *application) getWorkflowFile(filename string) dt.WorkflowFile {
	var workflowFile dt.WorkflowFile
	return getFile[dt.WorkflowFile](app, filename, workflowFile)
}
