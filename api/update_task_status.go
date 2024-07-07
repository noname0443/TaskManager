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
// @Param        status body updateTaskStatus true "Status"
// @Param        taskId path int true "TaskID"
// @Success      200  {string} string "ok"
// @Failure      400  {object} httputil.HTTPError
// @Failure      404  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /tasks/{taskId} [put]
func (c *Controller) UpdateTaskStatus(ctx *gin.Context) {
	req := updateTaskStatus{}
	if err := req.fromContext(ctx); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	fields := logrus.Fields{
		"taskId": req.taskId,
		"status": req.Status,
	}

	if err := c.updateTaskStatus(req.taskId, req.Status); err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debug("UpdateTaskStatus")
	ctx.String(200, "ok")
}

type updateTaskStatus struct {
	Status bool `json:"status"`
	taskId int
}

func (req *updateTaskStatus) fromContext(ctx *gin.Context) (err error) {
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return err
	}

	req.taskId, err = strconv.Atoi(ctx.Param("taskId"))
	if err != nil {
		return fmt.Errorf(httputil.INCORRECT_FORMAT, "taskId")
	}

	return nil
}

func (c *Controller) updateTaskStatus(taskId int, status bool) error {
	whereId := map[string]interface{}{"id": taskId}

	task := models.Task{}
	task.ID = uint(taskId)

	if err := c.db.Where(whereId).First(&task).Error; err != nil {
		return err
	}

	if !task.Status && status {
		if err := c.db.Model(&models.Task{}).Where(whereId).Updates(&models.Task{
			Status: true,
			Start:  time.Now(),
		}).Error; err != nil {
			return err
		}
	} else if task.Status && !status {
		if err := c.db.Model(&models.TimeSpent{}).Save(&models.TimeSpent{
			TaskID:        task.ID,
			BeginInterval: task.Start,
			EndInterval:   time.Now(),
		}).Error; err != nil {
			return err
		}

		task.Status = false
		task.Start = time.Time{}

		if err := c.db.Model(&models.Task{}).Where(whereId).Save(&task).Error; err != nil {
			return err
		}
	}

	return nil
}
