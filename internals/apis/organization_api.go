package apis

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type OrganizationApi struct {
	orgSvc *services.OrganizationService
}

func RegisterOrganizationApi(parentGroup *gin.RouterGroup, orgSvc *services.OrganizationService) {
	group := parentGroup.Group("/organizations")
	orgApi := &OrganizationApi{
		orgSvc: orgSvc,
	}

	group.GET("/", orgApi.GetListOfOrganizations)
	group.POST("/", orgApi.CreateOrganization)

	group.GET("/:orgId", orgApi.GetOrganizationFromId)
	group.PUT("/:orgId", orgApi.UpdateOrganization)
	group.DELETE("/:orgId")
}

func (orgApi *OrganizationApi) CreateOrganization(c *gin.Context) {
	orgRequest := NewOrganizationRequest{}

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&orgRequest); err != nil {
		responseError(c, http.StatusBadRequest, "Unable to process request body")
		return
	}

	exists := orgApi.orgSvc.ExistsByName(orgRequest.Name)
	if exists {
		responseError(c, http.StatusConflict, "Organization already exists")
		return
	}

	newOrg, err := orgApi.orgSvc.NewOrganization(orgRequest.Name)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong on our side")
		return
	}

	responseSingleItemStatus(c, http.StatusCreated, newOrg)
}

func (orgApi *OrganizationApi) GetListOfOrganizations(c *gin.Context) {

}

func (orgApi *OrganizationApi) GetOrganizationFromId(c *gin.Context) {
	orgId := c.Param("orgId")
	test := c.GetString("x-organization-id")
	log.Printf("Got context: %s\n", test)

	org, err := orgApi.orgSvc.GetOrganizationById(orgId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong our end")
		return
	}

	if org == nil {
		responseError(c, http.StatusNotFound, "No organization with that ID found")
		return
	}

	responseSingleItem(c, org)
}

func (orgApi *OrganizationApi) UpdateOrganization(c *gin.Context) {

}
