package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/maslow123/bookings/cmd/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"not-exist", "/omama/olala", "GET", http.StatusNotFound},
	// new routes
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"reservations-new", "/admin/reservations-new", "GET", http.StatusOK},
	{"reservations-all", "/admin/reservations-all", "GET", http.StatusOK},
	{"reservations-new with param", "/admin/reservations/new/1/show", "GET", http.StatusOK},

	// {"post-search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},

	// {"post-search-availability-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},

	// {"make reservation post", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "Omama"},
	// 	{key: "last_name", value: "Olala"},
	// 	{key: "email", value: "omama@getnada.com"},
	// 	{key: "phone", value: "878038246298102"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTest {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("For %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2020-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2020-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Omama")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Olala")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=omama@getnada.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=11111111")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2020-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Omama")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Olala")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=omama@getnada.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=11111111")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid end date
	reqBody = "start_date=2020-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Omama")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Olala")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=omama@getnada.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=11111111")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid roomID
	reqBody = "start_date=2020-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2020-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Omama")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Olala")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=omama@getnada.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=11111111")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=aaa")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code for invalid roomID: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid data
	reqBody = "start_date=2020-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2020-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=O")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Olala")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=omama@getnada.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=11111111")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert reservation into database
	reqBody = "start_date=2020-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2020-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Omama")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Olala")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=omama@getnada.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=11111111")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler failed when trying to fail insert reservation: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert room restriction into database
	reqBody = "start_date=2020-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2020-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Omama")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Olala")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=omama@getnada.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=11111111")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler failed when trying to fail insert room restriction: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// first case - rooms are not available

	postedData := url.Values{}
	postedData.Add("start", "2020-01-01")
	postedData.Add("end", "2020-01-02")
	postedData.Add("room_id", "1")

	// create the request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set http request header
	req.Header.Set("Content-Type", "x-www-form-url-encoded")

	// make handler handlerFunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// get response recorder
	rr := httptest.NewRecorder()

	// make request to our handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed to parse json")
	}

}

var loginTest = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"omama@getnada.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"boobee",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTest {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}

	}
}

// gets the context
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
