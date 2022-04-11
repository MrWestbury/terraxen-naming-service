package apis

import (
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type SchemaApi struct{}

func RegisterSchemaApi(parentGroup *gin.RouterGroup, svc *services.SchemaService) {
	group := parentGroup.Group("/schemas")

	schemaApi := &SchemaApi{}

	group.GET("/", schemaApi.NotImplemented)
	group.POST("/", schemaApi.NotImplemented)

	group.GET("/:schema", schemaApi.NotImplemented) // List schema versions

	group.GET("/:schema/:version", schemaApi.NotImplemented)
	group.PUT("/:schema/:version", schemaApi.NotImplemented)
	group.DELETE("/:schema/:version", schemaApi.NotImplemented)

	group.POST("/:schema/:version/resolve", schemaApi.ResolveResourceName)
}

func (sApi *SchemaApi) NotImplemented(c *gin.Context) {
	responseError(c, http.StatusNotImplemented, "Not yet implemented")
}

func (sApi *SchemaApi) ResolveResourceName(c *gin.Context) {

}
