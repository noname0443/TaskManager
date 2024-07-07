package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/noname0443/task_manager/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateTask godoc
//
// @Summary      Creates task
// @Description  Creates task
// @Accept       json
// @Produce      json
// @Param        request body CreateTaskReq true "CreateTaskReq"
// @Success      200  {int} taskId
// @Failure      400  {object} httputil.HTTPError
// @Failure      404  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /tasks/ [post]
func (c *Controller) CreateTask(ctx *gin.Context) {
	req := CreateTaskReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	fields := logrus.Fields{
		"description": req.Description,
		"userId":      req.UserID,
	}

	userExists, taskId, err := c.createTask(req.UserID, req.Description)
	if !userExists {
		logrus.WithFields(fields).Debug(err)
		httputil.NewError(ctx, 404, fmt.Errorf(httputil.NOT_FOUND))
		return
	}
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debug("CreateTask")
	ctx.JSON(200, taskId)
}

type CreateTaskReq struct {
	UserID      int    `json:"userId"`
	Description string `json:"description"`
}

func (c *Controller) createTask(userId int, description string) (userExists bool, taskId int, err error) {
	result := c.db.First(&models.User{}, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, 0, nil
	} else if result.Error != nil {
		return true, 0, result.Error
	}

	task := models.Task{
		UserID:      uint(userId),
		Description: description,
	}

	if err := c.db.Create(&task).Error; err != nil {
		return true, 0, err
	}

	return true, int(task.ID), nil
}
