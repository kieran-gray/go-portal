package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"html/template"
	S3Client "github.com/kieran-gray/go-portal/pkg/s3Client"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	templateCache map[string]*template.Template
	fileCache map[string]interface{}
	s3Client S3Client.S3
}

const templateRootDir string = "./ui/html/"

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	templates := []string{"index", "admin"}
	templateCache, err := templateCache(templateRootDir, templates)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		templateCache: templateCache,
		fileCache: map[string]interface{} {},
		s3Client: S3Client.S3Client(),
	}

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("---- Starting Server on %s ----", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}