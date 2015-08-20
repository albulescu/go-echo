package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type EndpointHandler func(request Request, response Response)

type Request struct {
	httpRequest *http.Request
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
	endpoints map[string]EndpointHandler
}

func (r *Router) Init() {
	r.endpoints = make(map[string]EndpointHandler)
}

func (r *Router) Map(uri string, method string, handler EndpointHandler) {
	key := fmt.Sprintf("%s-%s", method, uri)
	if _, ok := r.endpoints[key]; ok {
		panic("Route already registered")
	}
	r.endpoints[key] = handler
}

func (r *Router) Route(request *http.Request) (EndpointHandler, error) {

	key := fmt.Sprintf("%s-%s", request.Method, request.URL.Path)

	if handler, ok := r.endpoints[key]; ok {
		return handler, nil
	}

	return nil, &ErrorMessage{"Route does not exist"}
}

var router = new(Router)

/**
 * Call the route Endpoint
 */
func HttpApiHandle(rw http.ResponseWriter, hr *http.Request) {

	endpointHandler, err := router.Route(hr)

	if err != nil {
		http.Error(rw, "Route not found", 404)
		return
	}

	endpointHandler(Request{hr}, Response{rw})
}

/**
 * Http handler for api endpoint
 */
func ApiRouterInit() {

	router.Init()

	//map endpoints
	router.Map("/api/info", "GET", ApiInfo)
	router.Map("/api/auth", "POST", ApiAuth)

	//handle socket
	http.HandleFunc("/api/", HttpApiHandle)
	http.ListenAndServe(configApi.BindAddress, nil)
}
