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
	Route{"AddCallback", "POST", "/api/applications/{applicationId}/addcallback", applications.AddCallbackHandler},
	Route{"PostApplication", "POST", "/api/applications", applications.ApplicationAddHandler},
	Route{"UrlList", "GET", "/api/applications/{applicationId}/urls", applications.UrlListHandler},
	Route{"PostUrl", "POST", "/api/applications/{applicationId}/urls", applications.UrlAddHandler},
	Route{"FetchUrl", "GET", "/api/applications/{applicationId}/fetch/{urlId}", applications.FetchURLHandler},
	Route{"FetchUrls", "GET", "/api/applications/{applicationId}/fetch", applications.FetchApplicationURLs},
	Route{"UrlDeleteHandler", "DELETE", "/api/applications/{applicationId}/urls/{urlId}", applications.UrlDeleteHandler},
	Route{"UrlDetailHandler", "GET", "/api/applications/{applicationId}/urls/{urlId}", applications.UrlDetailHandler},
	Route{"UrlEditHandler", "PUT", "/api/applications/{applicationId}/urls/{urlId}", applications.UrlEditHandler},

	Route{"SignUp", "POST", "/api/signup", users.SignUpHandler},
	Route{"Login", "POST", "/api/login", users.LoginHandler},


	Route{"IndexHandler", "GET", "/", applications.IndexHandler},
	Route{"SignupTemplateHandler", "GET", "/signup", users.SignupTemplateHandler},
	Route{"LoginTemplateHandler", "GET", "/login", users.LoginTemplateHandler},
	Route{"ApplicationListTemplateHandler", "GET", "/applications", applications.ApplicationListTemplateHandler},
	Route{"UrlListTemplateHandler", "GET", "/applications/{applicationId}/urls", applications.UrlListTemplateHandler},
	Route{"UrlDetailTemplateHandler", "GET", "/applications/{applicationId}/urls/{urlId}", applications.UrlDetailTemplateHandler},
}
