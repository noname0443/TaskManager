package api

import (
	"errors"

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
// @Param        userId body uint true "User ID"
// @Param        description body string true "Description"
// @Success      200  {int} taskId
// @Failure      400  {object} httputil.HTTPError
// @Failure      404  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /tasks/ [post]
func (c *Controller) CreateTask(ctx *gin.Context) {
	req := CreateTaskRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	user := models.User{}
	user.ID = req.UserID

	result := c.db.First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logrus.Debug(req, result.Error)
		httputil.NewError(ctx, 404, result.Error)
		return
	} else if result.Error != nil {
		logrus.Debug(req, result.Error)
		httputil.NewError(ctx, 500, result.Error)
		return
	}

	task := models.Task{
		UserID:      req.UserID,
		Description: req.Description,
	}
	c.db.Create(&task)
	logrus.Debug("CreateTask ", task)
	ctx.JSON(200, task.ID)
}

type CreateTaskRequest struct {
	UserID      uint   `json:"userId"`
	Description string `json:"description"`
}
