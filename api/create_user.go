package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/noname0443/task_manager/integration"
	"github.com/sirupsen/logrus"
)

// CreateUser godoc
//
// @Summary      Creates user
// @Description  Creates user
// @Accept       json
// @Produce      json
// @Param        passportNumber body CreateUserReq true "Create user"
// @Success      200  {uint} userId
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/ [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	req := CreateUserReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	fields := logrus.Fields{
		"passportNumber": req.PassportNumber,
	}

	user, err := integration.GetPeopleInfo(req.PassportNumber)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, errors.New(httputil.SOMETHING_WENT_WRONG))
		return
	}

	if err := c.db.Create(&user).Error; err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, errors.New(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debug("CreateUser")
	ctx.JSON(200, user.ID)
}

type CreateUserReq struct {
	PassportNumber string `json:"passportNumber"`
}
