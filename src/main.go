package main

import (
	"fmt"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"

	"github.com/groundstorm/cobalt/src/api"
)

func serveSwaggerUI(c *restful.Container) {
	openAPIConfig := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/apidocs.json",
	}
	swaggerUIPath := "/apidocs/"
	swaggerDir := "./frontend/swagger-ui/dist"

	c.Add(restfulspec.NewOpenAPIService(openAPIConfig))
	c.ServeMux.Handle(swaggerUIPath, http.StripPrefix(swaggerUIPath, http.FileServer(http.Dir(swaggerDir))))
}

func main() {
	// Create the api object
	container := restful.NewContainer()
	api.CreateAPI(container)
	serveSwaggerUI(container)

	// Start the server
	fmt.Println("Starting server on http://localhost:1337/")
	server := &http.Server{Addr: ":1337", Handler: container}
	log.Fatal(server.ListenAndServe())
}
