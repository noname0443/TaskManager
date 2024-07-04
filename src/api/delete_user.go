package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
)

// DeleteUser godoc
//
// @Summary      Deletes user
// @Description  Deletes user
// @Accept       json
// @Produce      json
// @Param        userid path int true "User ID"
// @Success      200  {string}  string	"ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/{userid} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}
