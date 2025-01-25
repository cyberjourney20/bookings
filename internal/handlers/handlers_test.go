package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cyberjourney20/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var getTestVars = []struct {
	name              string
	url               string
	method            string
	expectedSatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

// TestHandlers tests all basic http "GET" request page handlers listed in getTestVars
func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range getTestVars {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedSatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedSatusCode, resp.StatusCode)
		}
	}
}

//	var testPAVars = []struct {
//		start_date        string
//		end_date          string
//		room_id           string
//		message           string
//		expectedSatusCode int
//	}{
//
//		{"2050-01-01", "2050-01-02", "2", "(all valid data)", http.StatusOK},
//		{"invalid", "2050-01-02", "2", "(invalid start_date)", http.StatusSeeOther},
//		{"2050-01-01", "invalid", "2", "(invalid end_date)", http.StatusSeeOther},
//		{"2050-01-01", "2050-01-02", "invalid", "(invalid room_id)", http.StatusSeeOther},
//	}
//
// TestPostAvailability tests the PostAvailabilityHandler
// testPostAvailabilityData is data for the PostAvailability handler test, /search-availability
var testPostAvailabilityData = []struct {
	name               string
	postedData         url.Values
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "rooms not available",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
		},
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name: "rooms are available",
		postedData: url.Values{
			"start_date": {"2040-01-01"},
			"end_date":   {"2040-01-02"},
			"room_id":    {"2"},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "empty post body",
		postedData:         url.Values{},
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name: "start date wrong format",
		postedData: url.Values{
			"start_date": {"invalid"},
			"end_date":   {"2040-01-02"},
			"room_id":    {"2"},
		},
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name: "end date wrong format",
		postedData: url.Values{
			"start_date": {"2040-01-01"},
			"end_date":   {"invalid"},
		},
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name: "database query fails",
		postedData: url.Values{
			"start_date": {"2060-01-01"},
			"end_date":   {"2060-01-02"},
		},
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestPostAvailability(t *testing.T) {
	for _, e := range testPostAvailabilityData {
		req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(e.postedData.Encode()))

		// get the context with session
		ctx := GetCtx(req)
		req = req.WithContext(ctx)

		// set the request header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// make our handler a http.HandlerFunc and call
		handler := http.HandlerFunc(Repo.PostAvailability)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s gave wrong status code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

// TestPostAvalibility tests the /search-availability "POST" handler
// func TestPostAvalibility(t *testing.T) {

// 	for _, e := range testPAVars {
// 		postedData := url.Values{}
// 		postedData.Add("start_date", e.start_date)
// 		postedData.Add("end_date", e.end_date)
// 		postedData.Add("room_id", e.room_id)
// 		rr, req, _ := POSTHandlersTestServer("POST", "/search-availability", postedData)
// 		// session.Put(ctx, "reservation", reservation)
// 		handler := http.HandlerFunc(Repo.SearchAvailability)
// 		handler.ServeHTTP(rr, req)
// 		if rr.Code != e.expectedSatusCode {
// 			t.Errorf("Reservation handler returned wrong response code for %s test: got %d, wanted %d", e.message, rr.Code, e.expectedSatusCode)
// 		}

// resp, err := ts.Client().POST(ts.URL + postedData)
// if err != nil {
// 	t.Log(err)
// 	t.Fatal(err)
// }
// if resp.StatusCode != e.expectedSatusCode {
// 	t.Errorf("for %s, expected %d but got %d", e.name, e.expectedSatusCode, resp.StatusCode)
// }
//}

// func POSTHandlersTestServer(gp, url string, reqBody url.Values) (rr *httptest.ResponseRecorder, req *http.Request, ctx context.Context) {
// 	req, _ = http.NewRequest(gp, url, strings.NewReader(reqBody.Encode()))
// 	ctx = GetCtx(req)
// 	req = req.WithContext(ctx)
// 	rr = httptest.NewRecorder()
// 	return rr, req, ctx

// user selects dates clicks ok
//dates good || dates malformed

// recast only cares about format, test good/bad
// Runs SearchAvalibilityForAllRooms in postgres.go // dont need to test here
// no avalilability redirect
// availability >> /choose-room redirect

// res := models.Reservation{
// 	StartDate: startDate,
// 	EndDate:   endDate,
// }
// m.App.Session.Put(r.Context(), "reservation", res)

// form.get dates error for malformed dates? start, end
// test DB for room availibility (use fake DB)
// check no avalibility error
// test render /choose-room page tested in TestReservation because it redirects there
//}

// TestAvalibilityJSON tests JSON POST for the /search-availability-json handler
func Test_AvalibilityJSON(t *testing.T) {
	// review JSON Test Lecture - 123
}

// add test for func book room redirect {"cr", "/choose-room/2", "GET", http.StatusOK} requires res.RoomName
// TestReservations tests the /reservations "GET" handler

// Test Cases -
// test case where reservation is not in session
// test with non existent room
// Dates?

func TestReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Generals Quarters",
		},
	}
	rr, req, ctx := GETHandlersTestServer("GET", "/reservation")
	// req, _ := http.NewRequest("GET", "/reservation", nil)
	// ctx := GetCtx(req)
	// req = req.WithContext(ctx)

	// rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non existent room
	req, _ = http.NewRequest("GET", "/reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestPostReservation tests the /reservation "Post" handler
func TestPostReservation(t *testing.T) {
	// rebuild this whole section to use an if statement

	// postedData := url.Values
	// postedData.Add("start_date=2050-01-01")
	// postedData.Add("end_date=2050-01-02")
	// postedData.Add("first_name=john")
	//req, _ := http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	//build a models.Reservation to pass through to satisfy my version of func PostReservation

	start_date, _ := reCastDates("2050-01-01")
	end_date, _ := reCastDates("2050-01-02")

	reservation := models.Reservation{
		StartDate: start_date,
		EndDate:   end_date,
		RoomID:    2,
		Room: models.Room{
			ID:       2,
			RoomName: "Generals Quarters",
		},
	}
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=john")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-5555-5555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ := http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	session.Put(ctx, "reservation", reservation)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for mission post body
	req, _ = http.NewRequest("POST", "/reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=john")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-5555-5555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid end date
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=john")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-5555-5555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for room id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=john")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-5555-5555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid room id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid form data
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=j")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-5555-5555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid form data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

// TestReservationSummary Tests the /reservation-summary GET page handler
