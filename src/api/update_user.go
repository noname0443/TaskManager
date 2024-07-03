package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/noname0443/task_manager/models"
)

type User struct {
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
// @Param        user body User true "User"
// @Success      200  {string}  string	"ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/ [put]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
	ctx.JSON(http.StatusOK, models.User{})
}
