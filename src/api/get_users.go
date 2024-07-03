package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
)

// GetUsers godoc
//
// @Summary      Gets users
// @Description  Gets users
// @Accept       json
// @Produce      json
// @Param        offset query int true "Pagination offset"
// @Param        limit query int true "Pagination limit"
// @Param        filters query []string false "Filters"
// @Success      200  {array}  User
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/ [get]
func (c *Controller) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}
