package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/priyankardasrpa/bookings/internal/config"
	"github.com/priyankardasrpa/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// What am I going to store in sessions
	gob.Register(models.Reservation{})

	// Change this to true when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWrite struct{}

func (mw *myWrite) Header() http.Header {
	var h http.Header
	return h
}

func (mw *myWrite) Write(b []byte) (int, error) {
	return len(b), nil
}

func (mw *myWrite) WriteHeader(statusCode int) {

}
