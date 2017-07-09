package helpers

import "github.com/gin-gonic/gin"

//ResponseJSON func
func ResponseJSON(c *gin.Context, status int, code int, message string, data interface{}) {

	if status >= 400 {
		c.JSON(status, gin.H{
			"code":    code,
			"message": message,
		})
	} else {
		c.JSON(status, gin.H{
			"code":    code,
			"data":    data,
			"message": message,
		})
	}
}

//ResponseSuccessJSON func
func ResponseSuccessJSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(200, gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	})
}

//ResponseNoContentJSON func
func ResponseNoContentJSON(c *gin.Context) {
	c.Status(204)
}

//ResponseEntityListJSON func
func ResponseEntityListJSON(c *gin.Context, code int, message string, entityList interface{}, metadata interface{}, total int) {
	c.JSON(200, gin.H{
		"code":  1,
		"data":  entityList,
		"total": total,
		//"metadata": metadata,
		"message": message,
	})
}

//ResponseCreatedJSON func
func ResponseCreatedJSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(201, gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	})
}

//ResponseAuthJSON func
func ResponseAuthJSON(c *gin.Context, code int, message string) {
	c.JSON(401, gin.H{
		"code":    code,
		"message": message,
	})
	c.Abort()
}

//ResponseNotFoundJSON func
func ResponseNotFoundJSON(c *gin.Context, code int, message string) {
	c.JSON(404, gin.H{
		"code":    code,
		"message": message,
	})
}

//ResponseBadRequestJSON func
func ResponseBadRequestJSON(c *gin.Context, code int, message interface{}) {
	c.JSON(400, gin.H{
		"code":    code,
		"message": message,
	})
	c.Abort()
}

//ResponseServerErrorJSON func
func ResponseServerErrorJSON(c *gin.Context) {
	c.Status(500)
	c.Abort()
}

//Response error API

//ResponseErrorsJSON func
func ResponseErrorsJSON(c *gin.Context, errors Errors) {
	c.JSON(400, errors)
	c.Abort()
}

//ResponseErrorJSON func
func ResponseErrorJSON(c *gin.Context, error ErrorDetail) {
	c.JSON(400, error)
	c.Abort()
}

//ResponseForbiddenJSON func
func ResponseForbiddenJSON(c *gin.Context, code int, message interface{}) {
	c.JSON(403, gin.H{
		"code":    code,
		"message": message,
	})
	c.Abort()
}
