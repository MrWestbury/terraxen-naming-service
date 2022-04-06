package apis

import (
	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type Api struct {
	router *gin.Engine
}

func NewApi(config *config.Config) *Api {
	api := &Api{
		router: gin.Default(),
	}

	apiGroup := api.router.Group("/api")

	v1Group := apiGroup.Group("/v1")

	orgService := services.NewOrganizationService()
	RegisterOrganizationApi(v1Group, orgService)

	return api
}

func (api *Api) Run(listener string) error {
	err := api.router.Run(listener)
	return err
}
