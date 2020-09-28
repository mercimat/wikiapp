package web

import (
    "net/http"
)

var handlerTestCases = []struct{
    description string
    mock        MockConnection
    req         string
    handler     string
    status      int
    content     string
}{
    {
        description:    "Request to / redirects to /view/frontPage",
        mock:           MockConnection{},
        req:            "/",
        handler:        "front",
        status:         http.StatusFound,
        content:        "<a href=\"/view/frontPage\">Found</a>.\n\n",
    },
    {
        description:    "Request to /view/page redirects to /edit/page if page is not known",
        mock:           MockConnection{ErrGetPage: true},
        req:            "/view/page",
        handler:        "view",
        status:         http.StatusFound,
        content:        "<a href=\"/edit/page\">Found</a>.\n\n",
    },
    {
        description:    "Request to /view/page displays the content of page",
        mock:           MockConnection{ErrGetPage: false},
        req:            "/view/page",
        handler:        "view",
        status:         http.StatusOK,
        content:        "<div class=\"content\">\n            <p>mocked page</p>\n        </div>\n",
    },
    {
        description:    "Request to /edit/page returns a blank edit box if page is not known",
        mock:           MockConnection{ErrGetPage: true},
        req:            "/edit/page",
        handler:        "edit",
        status:         http.StatusOK,
        content:        "<div><textarea name=\"body\" rows=\"20\" cols=\"80\"></textarea></div>",
    },
    {
        description:    "Request to /edit/page returns an edit box pre-filled with the body of page",
        mock:           MockConnection{ErrGetPage: false},
        req:            "/edit/page",
        handler:        "edit",
        status:         http.StatusOK,
        content:        "<div><textarea name=\"body\" rows=\"20\" cols=\"80\">mocked page</textarea></div>",
    },
}

var saveTestCases = []struct{
    description string
    mock        MockConnection
    req         string
    handler     string
    status      int
    content     string
}{
    {
        description:    "Request to /save/page redirects to /view/page if the database operation succeeded",
        mock:           MockConnection{ErrSavePage: false},
        req:            "/save/page",
        handler:        "save",
        status:         http.StatusFound,
        content:        "/view/page",
    },
    {
        description:    "Request to /save/page returns an error page if the database operation failed",
        mock:           MockConnection{ErrSavePage: true},
        req:            "/save/page",
        handler:        "save",
        status:         http.StatusInternalServerError,
        content:        "Failed to save page {\"page\" \"this is the body\"}\n",
    },
}

