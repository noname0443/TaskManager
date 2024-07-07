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

// UpdateUser godoc
//
// @Summary      Updates user
// @Description  Updates user
// @Accept       json
// @Produce      json
// @Param        user body updateUserReq true "User"
// @Success      200  {string}  string "ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/{userId} [put]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	req := updateUserReq{}
	if err := req.fromContext(ctx); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	fields := logrus.Fields{
		"passportNumber": req.PassportNumber,
		"surname":        req.Surname,
		"name":           req.Name,
		"patronymic":     req.Patronymic,
		"address":        req.Address,
		"userId":         req.userId,
	}

	user, err := updateUserReqToUserModel(&req)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	if err := c.updateUserDB(req.userId, user); err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debug("UpdateUser")
	ctx.String(200, "ok")
}

type updateUserReq struct {
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
	userId         int
}

func updateUserReqToUserModel(req *updateUserReq) (*models.User, error) {
	byteArray, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	user := models.User{}
	err = json.Unmarshal(byteArray, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Controller) updateUserDB(userId int, user *models.User) error {
	return c.db.Model(&models.User{}).Where(map[string]interface{}{"id": userId}).Updates(&user).Error
}

func (req *updateUserReq) fromContext(ctx *gin.Context) (err error) {
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return err
	}

	req.userId, err = strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		return fmt.Errorf(httputil.INCORRECT_FORMAT, "userId")
	}

	return nil
}
