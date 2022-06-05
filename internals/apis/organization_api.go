package apis

import (
	"encoding/json"
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgSvc services.OrganizationServiceProvider
}

func NewOrganizationHandler(orgSvc services.OrganizationServiceProvider) *OrganizationHandler {
	orgHandler := &OrganizationHandler{
		orgSvc: orgSvc,
	}
	return orgHandler
}

func (orgApi *OrganizationHandler) CreateOrganization(c *gin.Context) {
	orgRequest := NewOrganizationRequest{}

	if err := DecodeBody(c, &orgRequest); err != nil {
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

func (orgApi *OrganizationHandler) GetListOfOrganizations(c *gin.Context) {

}

func (orgApi *OrganizationHandler) GetOrganization(c *gin.Context) {
	orgUrlId := c.Param("orgId")
	orgId := c.GetString(ORG_CONTEXT_NAME)

	if orgUrlId != orgId {
		responseError(c, http.StatusInternalServerError, "Organization ID mismatch")
		return
	}

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

func (orgApi *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	orgUrlId := c.Param("orgId")
	orgId := c.GetString(ORG_CONTEXT_NAME)

	if orgUrlId != orgId {
		responseError(c, http.StatusInternalServerError, "Organization ID mismatch")
		return
	}

	org, err := orgApi.orgSvc.GetOrganizationById(orgId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong our end")
		return
	}

	if org == nil {
		responseError(c, http.StatusNotFound, "No organization with that ID found")
		return
	}

	var updateReq UpdateOrganizationRequest
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&updateReq); err != nil {
		responseError(c, http.StatusBadRequest, "Unable to process request body")
		return
	}

	orgApi.orgSvc.UpdateOrganization(orgId, updateReq.Name, updateReq.Variables)

	responseSingleItem(c, org)
}

func (orgApi *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	orgUrlId := c.Param("orgId")
	orgId := c.GetString(ORG_CONTEXT_NAME)

	if orgUrlId != orgId {
		responseError(c, http.StatusInternalServerError, "Organization ID mismatch")
		return
	}

	orgApi.orgSvc.DeleteOrganization(orgId)
}
