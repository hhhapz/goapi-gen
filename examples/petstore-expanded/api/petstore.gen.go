// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version (devel) DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/discord-gophers/goapi-gen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Error defines model for Error.
type Error struct {
	// Error message
	Message string `json:"message"`
}

// NewPet defines model for NewPet.
type NewPet struct {
	// Name of the pet
	Name string `json:"name"`

	// Type of the pet
	Tag *string `json:"tag,omitempty"`
}

// Pet defines model for Pet.
type Pet struct {
	// Embedded struct due to allOf(#/components/schemas/NewPet)
	NewPet `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// Unique id of the pet
	ID int64 `json:"id"`
}

// FindPetsParams defines parameters for FindPets.
type FindPetsParams struct {
	// tags to filter by
	Tags []string `json:"tags,omitempty"`

	// maximum number of results to return
	Limit *int32 `json:"limit,omitempty"`
}

// AddPetJSONBody defines parameters for AddPet.
type AddPetJSONBody NewPet

// AddPetJSONRequestBody defines body for AddPet for application/json ContentType.
type AddPetJSONRequestBody AddPetJSONBody

// Bind implements render.Binder.
func (AddPetJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
type Response struct {
	body        interface{}
	statusCode  int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.statusCode)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(statusCode int) *Response {
	resp.statusCode = statusCode
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// FindPetsJSON200Response is a constructor method for a FindPets response.
// A *Response is returned with the configured status code and content type from the spec.
func FindPetsJSON200Response(body []Pet) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// FindPetsJSONDefaultResponse is a constructor method for a FindPets response.
// A *Response is returned with the configured status code and content type from the spec.
func FindPetsJSONDefaultResponse(body Error) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// AddPetJSON201Response is a constructor method for a AddPet response.
// A *Response is returned with the configured status code and content type from the spec.
func AddPetJSON201Response(body Pet) *Response {
	return &Response{
		body:        body,
		statusCode:  201,
		contentType: "application/json",
	}
}

// AddPetJSONDefaultResponse is a constructor method for a AddPet response.
// A *Response is returned with the configured status code and content type from the spec.
func AddPetJSONDefaultResponse(body Error) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// DeletePetJSONDefaultResponse is a constructor method for a DeletePet response.
// A *Response is returned with the configured status code and content type from the spec.
func DeletePetJSONDefaultResponse(body Error) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// FindPetByIDJSON200Response is a constructor method for a FindPetByID response.
// A *Response is returned with the configured status code and content type from the spec.
func FindPetByIDJSON200Response(body Pet) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// FindPetByIDJSONDefaultResponse is a constructor method for a FindPetByID response.
// A *Response is returned with the configured status code and content type from the spec.
func FindPetByIDJSONDefaultResponse(body Error) *Response {
	return &Response{
		body:        body,
		statusCode:  200,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns all pets
	// (GET /pets)
	FindPets(w http.ResponseWriter, r *http.Request, params FindPetsParams)
	// Creates a new pet
	// (POST /pets)
	AddPet(w http.ResponseWriter, r *http.Request)
	// Deletes a pet by ID
	// (DELETE /pets/{id})
	DeletePet(w http.ResponseWriter, r *http.Request, id int64)
	// Returns a pet by ID
	// (GET /pets/{id})
	FindPetByID(w http.ResponseWriter, r *http.Request, id int64)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	Middlewares      map[string]func(http.Handler) http.Handler
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// FindPets operation middleware
func (siw *ServerInterfaceWrapper) FindPets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parameter object where we will unmarshal all parameters from the context
	var params FindPetsParams

	// ------------- Optional query parameter "tags" -------------

	if err := runtime.BindQueryParameter("form", true, false, "tags", r.URL.Query(), &params.Tags); err != nil {
		err = fmt.Errorf("invalid format for parameter tags: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	// ------------- Optional query parameter "limit" -------------

	if err := runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit); err != nil {
		err = fmt.Errorf("invalid format for parameter limit: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FindPets(w, r, params)
	})

	handler(w, r.WithContext(ctx))
}

// AddPet operation middleware
func (siw *ServerInterfaceWrapper) AddPet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddPet(w, r)
	})

	handler(w, r.WithContext(ctx))
}

// DeletePet operation middleware
func (siw *ServerInterfaceWrapper) DeletePet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id int64

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		err = fmt.Errorf("invalid format for parameter id: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeletePet(w, r, id)
	})

	handler(w, r.WithContext(ctx))
}

// FindPetByID operation middleware
func (siw *ServerInterfaceWrapper) FindPetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id int64

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		err = fmt.Errorf("invalid format for parameter id: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FindPetByID(w, r, id)
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	error
}
type UnmarshalingParamError struct {
	error
}
type RequiredParamError struct {
	error
}
type RequiredHeaderError struct {
	error
}
type InvalidParamFormatError struct {
	error
}
type TooManyValuesForParamError struct {
	error
}

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      map[string]func(http.Handler) http.Handler
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:     "/",
		BaseRouter:  chi.NewRouter(),
		Middlewares: make(map[string]func(http.Handler) http.Handler),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		Middlewares:      options.Middlewares,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Get("/pets", wrapper.FindPets)
		r.Post("/pets", wrapper.AddPet)
		r.Delete("/pets/{id}", wrapper.DeletePet)
		r.Get("/pets/{id}", wrapper.FindPetByID)

	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithMiddleware(key string, middleware func(http.Handler) http.Handler) ServerOption {
	return func(s *ServerOptions) {
		s.Middlewares[key] = middleware
	}
}

func WithMiddlewares(middlewares map[string]func(http.Handler) http.Handler) ServerOption {
	return func(s *ServerOptions) {
		s.Middlewares = middlewares
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"KAoSAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
