package main

import (
	"fmt"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	spec "github.com/go-openapi/spec"

	"github.com/groundstorm/cobalt/apps/go/src/cobalt/api"
)

func serveSwaggerUI(c *restful.Container) {
	addInfoToSwagger := func(s *spec.Swagger) {
		s.Info = &spec.Info{
			InfoProps: spec.InfoProps{
				Version:     "0.0.1",
				Title:       "Cobalt Event Server",
				Description: "cobalt event server",
			},
		}
	}

	openAPIConfig := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/apidocs.json",
		PostBuildSwaggerObjectHandler: addInfoToSwagger,
	}
	swaggerUIPath := "/apidocs/"
	swaggerDir := "./apps/go/static/swagger-ui/dist"

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
