package main

import (
	"encoding/json"
	"fmt"
	"github.com/dvsekhvalnov/jose2go"
	"log"
	"net/http"
	"strings"
)

const (
	INTERNAL_ERROR_CONFIG_INVALID = 1000
)

type JWTPayload struct {
	Id  string `json:"_id"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

type ResponseError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (re *ResponseError) Json() []byte {

	jsonResponse, err := json.Marshal(re)

	if err != nil {
		log.Print("Fail to encode json", err)
		return nil
	}

	return jsonResponse
}

type EndpointHandler func(request Request, response Response)

type Request struct {
	HttpRequest *http.Request
	user        map[string]interface{}
}

func (r *Request) IsLogged() bool {
	if r.user != nil {
		return true
	}
	return false
}

func (r *Request) SetUser(u map[string]interface{}) {
	r.user = u
}

func (r *Request) GetUser() map[string]interface{} {
	return r.user
}

/**
 * Response class that map the http response writter
 */

type Response struct {
	writter http.ResponseWriter
}

func (r *Response) Success(data interface{}, code int) {
	jsonResponse, err := json.Marshal(data)
	r.StatusCode(code)
	if err == nil {
		r.writter.Write(jsonResponse)
	} else {
		log.Print("Fail to encode json", err)
	}
}

func (r *Response) InternalError(code int) {
	r.Error("Internal Error", code, 500)
}

func (r *Response) StatusCode(code int) {
	r.writter.WriteHeader(code)
}

func (r *Response) Error(message string, code int, statusCode int) {

	errrorStruct := ResponseError{
		Message: message,
		Code:    code,
	}

	jsonResponse, err := json.Marshal(errrorStruct)

	if err != nil {
		log.Print("Fail to encode json", err)
		return
	}

	r.StatusCode(statusCode)
	r.writter.Write(jsonResponse)
}

/**
 * Router class used to keep and execute matched route
 */

type Router struct {
	routes map[string]Route
}

func (r *Router) Init() {
	r.routes = make(map[string]Route)
}

func (r *Router) Shutdown() {

}

func (r *Router) Map(uri string, route Route) {
	if _, ok := r.routes[uri]; ok {
		die("Route already registered")
	}
	r.routes[uri] = route
}

func (r *Router) Route(request *http.Request) (Route, error) {
	if route, ok := r.routes[request.URL.Path]; ok {
		return route, nil
	}

	return Route{}, &ErrorMessage{"Route does not exist"}
}

type Route struct {
	Method        string
	Authorization bool
	Handler       EndpointHandler
}

func CreateError(message string, code int) string {
	err := ResponseError{
		Message: message,
		Code:    code,
	}
	return string(err.Json())
}

var router = new(Router)

/**
 * Call the route Endpoint
 */
func HttpApiHandle(rw http.ResponseWriter, hr *http.Request) {

	route, err := router.Route(hr)

	if err != nil {
		http.Error(rw, CreateError("Route not found", 0), 404)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if route.Method != hr.Method {
		http.Error(rw, CreateError(fmt.Sprintf("Route accept only %s", route.Method), 0), 404)
		return
	}

	request := Request{}
	request.HttpRequest = hr

	response := Response{rw}

	if route.Authorization {

		authorization := hr.Header.Get("Authorization")

		if authorization == "" {
			http.Error(rw, CreateError("Unautorized", 0), 401)
			return
		}

		parts := strings.Split(authorization, "Bearer ")

		if len(parts) != 2 {
			http.Error(rw, CreateError("Invalid authorization header", 0), 400)
			return
		}

		bearer := parts[1]

		payload, _, err := jose.Decode(bearer, []byte(configApi.SecretKey))

		if err != nil {
			http.Error(rw, CreateError("Invalid bearer", 0), 401)
			return
		}

		jwt := JWTPayload{}

		err = json.Unmarshal([]byte(payload), &jwt)

		if err != nil {
			http.Error(rw, CreateError("Invalid JWT Format", 0), 500)
			return
		}

		log.Print("Load user with id:", jwt.Id)

		//TODO:
		// Load user with id jwt.Id
		// request.SetUser( user )

		request.SetUser(map[string]interface{}{
			"name": "Dummy",
		})
	}

	route.Handler(request, response)
}

/**
 * Http handler for api endpoint
 */
func ApiRouterInit() {

	router.Init()

	//map routes
	router.Map("/api/info", Route{"GET", false, ApiInfo})
	router.Map("/api/me", Route{"GET", true, ApiMe})
	router.Map("/api/auth", Route{"POST", false, ApiAuth})

	apiLog.Printf("Api started on %s", configApi.BindAddress)

	//handle socket
	http.HandleFunc("/api/", HttpApiHandle)
	http.ListenAndServe(configApi.BindAddress, nil)
}
