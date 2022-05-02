package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListMeta struct {
	Offset int
	Limit  int
}

func DecodeBody(c *gin.Context, result interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&result); err != nil {
		log.Printf("failed to decode body: %v", err)
		responseError(c, http.StatusBadRequest, "Unable to process request body")
		return err
	}
	return nil
}

func ProcessListMetadata(c *gin.Context) *ListMeta {
	lm := &ListMeta{
		Offset: 0,
		Limit:  50,
	}

	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err == nil {
		lm.Offset = offset
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err == nil {
		lm.Limit = limit
	}

	return lm
}

// Temp
func NotImplemented(c *gin.Context) {
	responseError(c, http.StatusNotImplemented, "Not yet implemented")
}
