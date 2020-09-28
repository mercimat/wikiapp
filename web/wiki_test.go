package web

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"
)

type MockConnection struct {
    ErrGetPage  bool
    ErrSavePage bool
}

func (m MockConnection) SavePage(p *Page) error {
    if m.ErrSavePage {
        return fmt.Errorf("Failed to save page %q", *p)
    }
    return nil
}

func (m MockConnection) GetPage(title string) (*Page, error) {
    if m.ErrGetPage {
        return &Page{}, fmt.Errorf("Failed to get page %s", title)
    }
    return &Page{title, []byte("mocked page")}, nil
}

func init() {
    InitTemplates("../tmpl/edit.html", "../tmpl/view.html")
}

func TestHandlers(t *testing.T) {
    for _, tc := range handlerTestCases {
        t.Logf("Test case: %s", tc.description)

        req, err := http.NewRequest("GET", tc.req, nil)
        if err != nil {
            t.Fatal(err)
        }

        var handler http.HandlerFunc
        switch tc.handler {
        case "edit":
            handler = http.HandlerFunc(MakeHandler(EditHandler, tc.mock))
        case "view":
            handler = http.HandlerFunc(MakeHandler(ViewHandler, tc.mock))
        default:
            handler = http.HandlerFunc(FrontPageHandler)
        }

        rr := httptest.NewRecorder()
        handler.ServeHTTP(rr, req)

        if rr.Code != tc.status {
            t.Fatalf("handler returned wrong status code: got %v want %v",
                rr.Code, tc.status)
        }

        if body := rr.Body.String(); !strings.Contains(body, tc.content) {
            t.Fatalf("handler returned wrong body: \nbody does not contain: %s\nGot %s",
                tc.content, body)
        }
        t.Logf("PASS: %s", tc.description)
    }
}

func TestSaveHandler(t *testing.T) {
    for _, tc := range saveTestCases {
        t.Logf("Test case: %s", tc.description)

        form := url.Values{}
        form.Add("body", "this is the body")

        req, err := http.NewRequest("POST", tc.req, strings.NewReader(form.Encode()))
        if err != nil {
            t.Fatal(err)
        }
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(MakeHandler(SaveHandler, tc.mock))
        handler.ServeHTTP(rr, req)

        if rr.Code != tc.status {
            t.Fatalf("handler returned wrong status code: got %v want %v",
                rr.Code, tc.status)
        }

        expected := tc.content

        if tc.status == http.StatusFound {
            // Body is empty in this case
            expected = ""

            // Verify that the user is redirected to /view/page
            loc, ok := rr.HeaderMap["Location"]
            if !ok {
                t.Fatalf("response header is missing Location information: %q", rr.HeaderMap)
            }
            if string(loc[0]) != tc.content {
                t.Fatalf("handler returned wrong location: got %s want %s", loc, tc.content)
            }
        }

        if body := rr.Body.String(); body != expected {
            t.Fatalf("handler returned wrong body: \nbody does not contain: %s\nGot %s",
                expected, body)
        }

        t.Logf("PASS: %s", tc.description)
    }
}

