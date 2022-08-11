package handlers

import (
	"net/http"

	"example.com/hello-world/pkg/config"
	"example.com/hello-world/pkg/models"
	"example.com/hello-world/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Store the client's ip address
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// Perform business logic
	stringMap := map[string]string{}
	stringMap["test"] = "Priyankar"

	// Retrive client's ip address from session
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	// Pass the remoteIP to template
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
