package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestFrontPageHandler(t *testing.T) {
    // Request to / redirects to the front page
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(frontPageHandler)

    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            rr.Code, http.StatusFound)
    }

    expectedBody := "<a href=\"/view/frontPage\">Found</a>.\n\n"
    if body := rr.Body.String(); body != expectedBody {
        t.Errorf("handler returned wrong body: got %s want %s",
            body, expectedBody)
    }
}

func TestViewHandler(t *testing.T) {
    // Request to /view/page redirects to /edit/page if the page is not known
    req, err := http.NewRequest("GET", "/view/page", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(makeHandler(viewHandler))

    handler.ServeHTTP(rr, req)

    t.Logf("Code: %d", rr.Code)
    t.Logf("Body: %q", rr.Body.String())

    //if rr.Code != http.StatusFound {
    //    t.Errorf("handler returned wrong status code: got %v want %v",
    //        rr.Code, http.StatusFound)
    //}

    //expectedBody := "<a href=\"/view/frontPage\">Found</a>.\n\n"
    //if body := rr.Body.String(); body != expectedBody {
    //    t.Errorf("handler returned wrong body: got %s want %s",
    //        body, expectedBody)
    //}
}
