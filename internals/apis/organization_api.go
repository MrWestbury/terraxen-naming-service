package apis

import (
	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type OrganizationApi struct {
	orgSvc *services.OrganizationService
}

func RegisterOrganizationApi(parentGroup *gin.RouterGroup, orgSvc *services.OrganizationService) {
	group := parentGroup.Group("/organization")
	orgApi := &OrganizationApi{
		orgSvc: orgSvc,
	}

	group.POST("/", orgApi.CreateOrganization)
}

func (orgApi *OrganizationApi) CreateOrganization(c *gin.Context) {

}
