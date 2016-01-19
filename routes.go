package main

import (
	"github.com/hipo/gotcha/applications"
	"github.com/hipo/gotcha/users"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{},
	Route{"ApplicationList", "GET", "/api/applications", applications.ApplicationListHandler},
	Route{"ApplicationDelete", "DELETE", "/api/applications/{applicationId}", applications.ApplicationDeleteHandler},

	Route{"PostApplication", "POST", "/api/applications", applications.ApplicationAddHandler},
	Route{"UrlList", "GET", "/api/applications/{applicationId}/urls", applications.UrlListHandler},
	Route{"PostUrl", "POST", "/api/applications/{applicationId}/urls", applications.UrlAddHandler},
	Route{"FetchUrls", "GET", "/api/applications/{applicationId}/fetch", applications.FetchApplicationURLs},
	Route{"UrlDeleteHandler", "DELETE", "/api/applications/{applicationId}/urls/{urlId}", applications.UrlDeleteHandler},

	Route{"SignUp", "POST", "/api/signup", users.SignUpHandler},
	Route{"Login", "POST", "/api/login", users.LoginHandler},

	Route{"SignupTemplateHandler", "GET", "/signup", users.SignupTemplateHandler},
	Route{"LoginTemplateHandler", "GET", "/login", users.LoginTemplateHandler},
	Route{"ApplicationListTemplateHandler", "GET", "/applications", applications.ApplicationListTemplateHandler},
	Route{"UrlListTemplateHandler", "GET", "/applications/{applicationId}/urls", applications.UrlListTemplateHandler},
}
