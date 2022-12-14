package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/priyankardasrpa/bookings/internal/config"
	"github.com/priyankardasrpa/bookings/internal/models"
)

var app *config.AppConfig
var pathToTemplates = "./templates"
var functions = template.FuncMap{}

func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default data to template data
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// Create the templates cache
	//tc, err := CreateTemplateCache()
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get the required template
	t, ok := tc[tmpl]
	if !ok {
		log.Println("Could not load template")
		return errors.New("can't get template from cache")
	}

	// More pinpoint error checking

	// Adding default template data
	td = AddDefaultData(td, r)

	buf := new(bytes.Buffer)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// Render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the file paths "./templates/*.page.tmpl"
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		fmt.Println(err)
		return myCache, err
	}

	// For each page
	for _, page := range pages {
		// Get the name of the page excluding the full path
		name := filepath.Base(page)

		// Get the parsed template
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println(err)
			return myCache, err
		}

		// Get all the *.layout.tmpl paths
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			fmt.Println(err)
			return myCache, err
		}

		if len(matches) > 0 {
			// Has at atleast 1 layout file to associate
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				fmt.Println(err)
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
