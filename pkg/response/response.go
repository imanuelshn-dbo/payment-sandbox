package response

import "github.com/gin-gonic/gin"

type Response struct {
	StatusCode    int         `json:"status_code"`
	SystemMessage string      `json:"system_message"`
	Data          interface{} `json:"data,omitempty"`
	Error         interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		StatusCode:    200,
		SystemMessage: "success",
		Data:          data,
	})
}

func Error(c *gin.Context, code int, err interface{}) {
	c.JSON(code, Response{
		StatusCode:    code,
		SystemMessage: "error",
		Error:         err,
	})
}
