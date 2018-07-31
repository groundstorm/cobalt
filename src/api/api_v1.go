package api

import (
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func CreateAPI(c *restful.Container) {
	createUsersAPI(c)
}

func createUsersAPI(c *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/api/v1")
	tags := []string{"users"}
	restful.Add(ws)

	ws.Route(
		ws.POST("/users").
			To(createUser).
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Reads(NewUser{}).
			Returns(201, "Created", OK{}).
			Returns(400, "Bad Request", Error{}))

	c.Add(ws)
}

func createUser(request *restful.Request, response *restful.Response) {
}
