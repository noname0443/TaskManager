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
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, fmt.Errorf(httputil.INCORRECT_FORMAT, "userId"))
		return
	}

	user := models.User{}
	user.ID = uint(userId)

	var exists bool
	if err = c.db.Model(&models.User{}).Select("count(*) > 0").Where(&user).Find(&exists).Error; err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	if !exists {
		logrus.Debug(err)
		httputil.NewError(ctx, 404, fmt.Errorf(httputil.NOT_FOUND))
		return
	}

	if err := c.db.Unscoped().Delete(&user).Error; err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}
	logrus.Debug("DeleteUser", userId)
	ctx.String(http.StatusOK, "ok")
}
