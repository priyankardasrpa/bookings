package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/priyankardasrpa/bookings/pkg/config"
	"github.com/priyankardasrpa/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default data to template data
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
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
		log.Fatal("Could not load template")
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
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the file paths "./templates/*.page.tmpl"
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		fmt.Println(err)
		return myCache, err
	}

	// For each page
	for _, page := range pages {
		// Get the name of the page excluding the full path
		name := filepath.Base(page)

		// Get the parsed template
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			fmt.Println(err)
			return myCache, err
		}

		// Get all the *.layout.tmpl paths
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			fmt.Println(err)
			return myCache, err
		}

		if len(matches) > 0 {
			// Has at atleast 1 layout file to associate
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				fmt.Println(err)
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
