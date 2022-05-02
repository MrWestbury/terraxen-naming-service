package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func responseError(c *gin.Context, statusCode int, errMessage string) {
	body := map[string]interface{}{
		"code":    statusCode,
		"message": errMessage,
	}

	c.AbortWithStatusJSON(statusCode, body)
}

func responseSingleItem(c *gin.Context, item interface{}) {
	responseSingleItemStatus(c, http.StatusOK, item)
}

func responseSingleItemStatus(c *gin.Context, statusCode int, item interface{}) {
	body := map[string]interface{}{
		"links": map[string]string{},
		"data":  item,
	}

	c.IndentedJSON(statusCode, body)
}

func responseNoContent(c *gin.Context, statusCode int) {
	c.Status(statusCode)
}
