package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
)

// GetUserTasks godoc
//
// @Summary      Gets the user's tasks
// @Description  Gets the user's tasks
// @Accept       json
// @Produce      json
// @Param        from query string true "from"
// @Param        to query string true "to"
// @Success      200  {array} UserTask
// @Failure      400  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /users/{userid} [get]
func (c *Controller) GetUserTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}
