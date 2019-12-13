/*
 * Otto user service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package internal

import (
	"github.com/danielmunro/otto-user-service/internal/controller"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/",
		Index,
	},

	{
		"CreateGroupV1",
		strings.ToUpper("Post"),
		"/group",
		CreateGroupV1,
	},

	{
		"CreateNewUser",
		strings.ToUpper("Post"),
		"/user",
		controller.CreateNewUserV1,
	},

	{
		"GetGroupV1",
		strings.ToUpper("Get"),
		"/group",
		GetGroupV1,
	},

	{
		"GetUserV1",
		strings.ToUpper("Get"),
		"/user/{uuid}",
		controller.GetUserV1,
	},

	{
		"CreateNewSession",
		strings.ToUpper("Post"),
		"/session",
		controller.CreateSessionV1,
	},

	{
		"RespondToChallenge",
		strings.ToUpper("Put"),
		"/session",
		controller.RespondToChallengeV1,
	},

	{
		"GetSession",
		strings.ToUpper("Get"),
		"/session",
		controller.GetSessionV1,
	},

	{
		"RefreshSession",
		strings.ToUpper("Patch"),
		"/session",
		controller.RefreshSessionV1,
	},
}
