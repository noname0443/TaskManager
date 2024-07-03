package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
)

// UpdateTaskStatus godoc
//
// @Summary      Update the task's status
// @Description  Update the task's status
// @Accept       json
// @Produce      json
// @Param        status query bool true "Status"
// @Param        taskid path int true "TaskID"
// @Success      200  {string} string "ok"
// @Failure      400  {object} httputil.HTTPError
// @Failure      404  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /tasks/{taskid} [put]
func (c *Controller) UpdateTaskStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}
