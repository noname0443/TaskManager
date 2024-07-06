package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/noname0443/task_manager/models"
	"github.com/sirupsen/logrus"
)

// UpdateTaskStatus godoc
//
// @Summary      Update the task's status
// @Description  Update the task's status
// @Accept       json
// @Produce      json
// @Param        status body UpdateTaskStatus true "Status"
// @Param        taskId path int true "TaskID"
// @Success      200  {string} string "ok"
// @Failure      400  {object} httputil.HTTPError
// @Failure      404  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /tasks/{taskId} [put]
func (c *Controller) UpdateTaskStatus(ctx *gin.Context) {
	taskId, err := strconv.Atoi(ctx.Param("taskId"))
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, fmt.Errorf(httputil.INCORRECT_FORMAT, "taskId"))
		return
	}

	req := UpdateTaskStatus{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	fields := logrus.Fields{
		"taskId": taskId,
		"status": req.Status,
	}

	task := models.Task{}
	task.ID = uint(taskId)

	if err := c.db.Where(&task).First(&task).Error; err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	if !task.Status && req.Status {
		if err := c.db.Model(&models.Task{}).Where(map[string]interface{}{"id": taskId}).Updates(&models.Task{
			Status: true,
			Start:  time.Now(),
		}).Error; err != nil {
			logrus.WithFields(fields).Warn(err)
			httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
			return
		}
	} else if task.Status && !req.Status {
		if err := c.db.Model(&models.TimeSpent{}).Save(&models.TimeSpent{
			TaskID:        task.ID,
			BeginInterval: task.Start,
			EndInterval:   time.Now(),
		}).Error; err != nil {
			logrus.WithFields(fields).Warn(err)
			httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
			return
		}

		task.Status = false
		task.Start = time.Time{}

		if err := c.db.Model(&models.Task{}).Where(map[string]interface{}{"id": taskId}).Save(&task).Error; err != nil {
			logrus.WithFields(fields).Warn(err)
			httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
			return
		}
	}
	logrus.WithFields(fields).Debug("UpdateTaskStatus")
	ctx.String(200, "ok")
}

type UpdateTaskStatus struct {
	Status bool `json:"status"`
}
