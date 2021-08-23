package render

import (
	"net/http"
	"testing"

	"github.com/maslow123/bookings/cmd/internal/config"
)

func TestAddDefaultData(t *testing.T) {
	var td config.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "xxx")
	result := AddDefaultData(&td, r)

	if result.Flash != "xxx" {
		t.Error("Flash value of xxx not found in session")
	}

}

func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	r = r.WithContext(ctx)

	return r, nil
}
