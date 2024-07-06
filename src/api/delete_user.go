package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/noname0443/task_manager/models"
	"github.com/sirupsen/logrus"
)

// DeleteUser godoc
//
// @Summary      Deletes user
// @Description  Deletes user
// @Accept       json
// @Produce      json
// @Param        userId path int true "UserID"
// @Success      200  {string}  string	"ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/{userId} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	req := deleteUserReq{}
	if err := req.fromContext(ctx); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	fields := logrus.Fields{
		"userId": req.userId,
	}

	exist, err := c.deleteUser(req.userId)
	if !exist {
		logrus.WithFields(fields).Debug(err)
		httputil.NewError(ctx, 404, fmt.Errorf(httputil.NOT_FOUND))
		return
	}
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debug("DeleteUser")
	ctx.String(http.StatusOK, "ok")
}

type deleteUserReq struct {
	userId int
}

func (req *deleteUserReq) fromContext(ctx *gin.Context) (err error) {
	req.userId, err = strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		return fmt.Errorf(httputil.INCORRECT_FORMAT, "userId")
	}

	return nil
}

func (c *Controller) deleteUser(userId int) (exists bool, err error) {
	result := c.db.Unscoped().Delete(&models.User{}, userId)
	if result.Error != nil {
		return true, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
