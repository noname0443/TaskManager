package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/noname0443/task_manager/models"
	"github.com/sirupsen/logrus"
)

type UpdateUserReq struct {
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

// UpdateUser godoc
//
// @Summary      Updates user
// @Description  Updates user
// @Accept       json
// @Produce      json
// @Param        user body UpdateUserReq true "User"
// @Success      200  {string}  string "ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/{userId} [put]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	req := UpdateUserReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, fmt.Errorf(httputil.INCORRECT_FORMAT, "userId"))
		return
	}

	fields := logrus.Fields{
		"UpdateUserReq": req,
		"userId":        userId,
	}

	byteArray, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	user := models.User{}
	err = json.Unmarshal(byteArray, &user)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	if err := c.db.Model(&models.User{}).Where(map[string]interface{}{"id": userId}).Updates(&user).Error; err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debug("UpdateUser")
	ctx.String(200, "ok")
}
