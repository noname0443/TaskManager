package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
)

// CreateTask godoc
//
// @Summary      Creates task
// @Description  Creates task
// @Accept       json
// @Produce      json
// @Param        description body string true "Description"
// @Success      200  {int} taskId
// @Failure      400  {object} httputil.HTTPError
// @Failure      404  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /tasks/ [post]
func (c *Controller) CreateTask(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}
