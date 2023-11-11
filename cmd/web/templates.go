package main

import (
	"html/template"
	"os"
	"path/filepath"
)

func getTemplateNames(dir string) []string {
	templateNames := []string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}
	for _, entry := range entries {
		if entry.IsDir() {
			templateNames = append(templateNames, entry.Name())
		}
	}
	return templateNames
}

func buildTemplate(path string, name string) (*template.Template, error) {
	template := template.New(name)
	functions, ok := templateToFuncMap[name]
	if ok {
		template = template.Funcs(functions)
	}

	dir := path + name + "/"
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		template, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		template, err = template.ParseGlob(filepath.Join(dir, "*.partial.tmpl")) //nolint:all
		if err != nil {
			return nil, err
		}
	}
	return template, nil
}

func templateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	for _, templateName := range getTemplateNames(dir) {
		template, err := buildTemplate(dir, templateName)
		if err != nil {
			return nil, err
		}
		cache[templateName] = template
	}
	return cache, nil
}
