package httputil

import "github.com/gin-gonic/gin"

const (
	SOMETHING_WENT_WRONG = "something went wrong"
	NOT_ENOUGHT_FIELDS   = "request does not have fields [%s]"
)

func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
