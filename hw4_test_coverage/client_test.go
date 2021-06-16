package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// код писать тут

type Userspec struct {
	ID            int    `xml:"id"`
	Guid          string `xml:"guid"`
	Active        bool   `xml:"isActive"`
	Balance       string `xml:"balance"`
	Picture       string `xml:"picture"`
	Age           int    `xml:"age"`
	EyeColor      string `xml:"eyeColor"`
	First_name    string `xml:"first_name"`
	Last_name     string `xml:"last_name"`
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favoriteFruit"`
}

type Roots struct {
	List []Userspec `xml:"row"`
}

type Usersspec struct {
	Version string `xml:"version,attr"`
	Rooter  Roots  `xml:"root"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("AccessToken")
	if token != "token123420" {
		w.WriteHeader(http.StatusUnauthorized)
	}
	order_field := r.FormValue("order_field")
	if order_field != "" && order_field != "Id" && order_field != "Age" && order_field != "Name" {
		w.WriteHeader(http.StatusBadRequest)
		if order_field == "testbadjson" {
			io.WriteString(w, `{"hahahahh`)
		} else if order_field == "testbadfieldnojson" {
			data, _ := json.Marshal(SearchErrorResponse{
				Error: "Undefind error",
			})
			w.Write(data)
		} else {
			data, _ := json.Marshal(SearchErrorResponse{
				Error: "ErrorBadOrderField",
			})
			w.Write(data)
		}
	} else {
		query := r.FormValue("query")
		if query == "BadUnpack" {
			io.WriteString(w, `{"hahahahh`)
		} else if query == "ServerError" {
			w.WriteHeader(http.StatusInternalServerError)
		} else if query == "Time" {
			time.Sleep(time.Second * 2)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `[{"Id": 42, "Name": "rvasily", "Age": 12, "About":"ok","Gender":"male"}]`)
		}
	}
}

func TestAccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "token",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      1,
		Offset:     1,
		Query:      "str",
		OrderField: "name",
		OrderBy:    -1,
	})
	if err == nil {
		t.Errorf("AccessToken does not work")
	}
}

func TestUrl(t *testing.T) {
	c := &SearchClient{
		URL:         "invalidurl",
		AccessToken: "token",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      1,
		Offset:     1,
		Query:      "str",
		OrderField: "name",
		OrderBy:    -1,
	})
	if err == nil {
		t.Errorf("Invalid URL")
	}
}

func TestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "token123420",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      1,
		Offset:     1,
		Query:      "Time",
		OrderField: "Name",
		OrderBy:    -1,
	})
	fmt.Print("sualalfaf")
	if err == nil {
		t.Errorf("Invalid URL")
	}
}

func TestOffset26(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "token123420",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      26,
		Offset:     1,
		Query:      "str",
		OrderField: "Name",
		OrderBy:    -1,
	})
	if err != nil {
		t.Errorf("There no error")
	}
}

func TestBasic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "token123420",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      0,
		Offset:     1,
		Query:      "str",
		OrderField: "Name",
		OrderBy:    -1,
	})
	if err != nil {
		t.Errorf("There no error")
	}
}

func TestUnpackBad(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "token123420",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      0,
		Offset:     1,
		Query:      "BadUnpack",
		OrderField: "Name",
		OrderBy:    -1,
	})
	if err == nil {
		t.Errorf("There no error")
	}
}

func Test(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "token123420",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      0,
		Offset:     1,
		Query:      "ServerError",
		OrderField: "Name",
		OrderBy:    -1,
	})
	if err == nil {
		t.Errorf("There no error")
	}
}

func TestLimitOffset(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "token",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      -1,
		Offset:     1,
		Query:      "str",
		OrderField: "name",
		OrderBy:    -1,
	})
	if err == nil {
		t.Errorf("Limit is less then 0")
	}
	_, err = c.FindUsers(SearchRequest{
		Limit:      25,
		Offset:     -1,
		Query:      "str",
		OrderField: "name",
		OrderBy:    -1,
	})
	if err == nil {
		t.Errorf("Offset is less then 0")
	}
}

func TestInvalidOrderField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		URL:         ts.URL,
		AccessToken: "token123420",
	}
	_, err := c.FindUsers(SearchRequest{
		Limit:      1,
		Offset:     1,
		Query:      "str",
		OrderField: "badfield",
		OrderBy:    1,
	})
	if err == nil {
		t.Errorf("OrderFeld is invalid")
	}
	_, err = c.FindUsers(SearchRequest{
		Limit:      1,
		Offset:     1,
		Query:      "str",
		OrderField: "testbadjson",
		OrderBy:    1,
	})
	if err == nil {
		t.Errorf("json is invalid")
	}
	_, err = c.FindUsers(SearchRequest{
		Limit:      1,
		Offset:     1,
		Query:      "str",
		OrderField: "testbadfieldnojson",
		OrderBy:    1,
	})
	if err == nil {
		t.Errorf("unknown bad request")
	}
}
