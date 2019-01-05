package main

import (
	"github.com/groundstorm/cobalt/user"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

type Services struct {
	users *user.UserService
}

func CreateWebService(s Services) *restful.WebService {	
	ws := &restful.WebService{}
	ws.Route(ws.GET("/hello").To(hello).
		Doc("create example thing").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(ExampleResp{}).
		Reads(ExampleReq{}).
		Returns(200, "OK", &ExampleResp{}).
	return ws
}

func hello(*restful.Request, *restful.Response) {	
}