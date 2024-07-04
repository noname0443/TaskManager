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

// GetUserTasks godoc
//
// @Summary      Gets the user's tasks
// @Description  Gets the user's tasks
// @Accept       json
// @Produce      json
// @Param        offset query int true "Pagination offset"
// @Param        limit query int true "Pagination limit"
// @Param        userId path int true "UserID"
// @Success      200  {array} UserTask
// @Failure      400  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /users/{userId} [get]
func (c *Controller) GetUserTasks(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, fmt.Errorf(httputil.INCORRECT_FORMAT, "limit"))
		return
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, fmt.Errorf(httputil.INCORRECT_FORMAT, "offset"))
		return
	}

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, fmt.Errorf(httputil.INCORRECT_FORMAT, "userId"))
		return
	}

	fields := logrus.Fields{
		"limit":  limit,
		"offset": offset,
		"userId": userId,
	}

	tasks := []models.Task{}
	if err := c.db.Where(map[string]interface{}{"userId": userId}).Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	tasksJson := make([]TaskJSON, len(tasks))
	for i, task := range tasks {
		tasksJson[i].UserID = task.UserID
		tasksJson[i].ID = task.ID
		tasksJson[i].Description = task.Description
		timeSpent := task.Duration
		if !task.Start.IsZero() {
			timeSpent += time.Since(task.Start)
		}
		tasksJson[i].Hours = uint(timeSpent.Hours())
		tasksJson[i].Minutes = uint(timeSpent.Minutes())
		tasksJson[i].Seconds = uint(timeSpent.Seconds())
		tasksJson[i].Status = task.Status
	}
	logrus.WithFields(fields).Debug("GetUserTasks")
	ctx.JSON(200, tasksJson)
}

type TaskJSON struct {
	UserID      uint   `json:"userId"`
	ID          uint   `json:"taskId"`
	Description string `json:"description"`
	Hours       uint   `json:"hours"`
	Minutes     uint   `json:"minutes"`
	Seconds     uint   `json:"seconds"`
	Status      bool   `json:"status"`
}
