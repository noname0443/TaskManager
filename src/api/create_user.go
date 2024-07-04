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
// @Param        passportNumber body CreateUserRequest true "Create user"
// @Success      200  {uint} userId
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/ [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	userData := CreateUserRequest{}
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	user, err := integration.GetPeopleInfo(userData.PassportNumber)
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 500, errors.New(httputil.SOMETHING_WENT_WRONG))
		return
	}

	c.db.Create(&user)
	logrus.Debug("CreateUser ", user)
	ctx.JSON(200, user.ID)
}

type CreateUserRequest struct {
	PassportNumber string `json:"passportNumber"`
}
