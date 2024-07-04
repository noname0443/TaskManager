package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/noname0443/task_manager/models"
	"github.com/sirupsen/logrus"
)

// GetUsers godoc
//
// @Summary      Gets users
// @Description  Gets users
// @Accept       json
// @Produce      json
// @Param        offset query int true "Pagination offset"
// @Param        limit query int true "Pagination limit"
// @Param        filters query []string false "Filters"
// @Success      200  {array} UserJSON
// @Failure      400  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /users/ [get]
func (c *Controller) GetUsers(ctx *gin.Context) {
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

	filters, err := ParseFilters(ctx.Query("filters"))
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	fields := logrus.Fields{
		"limit":   limit,
		"offset":  offset,
		"filters": filters,
	}

	users := []models.User{}
	if err := c.db.Where(filters).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	byteArray, err := json.Marshal(users)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	usersJson := []UserJSON{}
	err = json.Unmarshal(byteArray, &usersJson)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debugln("GetUsers")
	ctx.JSON(http.StatusOK, usersJson)
}

type UserJSON struct {
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

func ParseFilters(filtersRaw string) (map[string]string, error) {
	if len(filtersRaw) == 0 {
		return nil, nil
	}
	filtersPairs := strings.Split(filtersRaw, ",")
	filters := map[string]string{}
	for _, filter := range filtersPairs {
		values := strings.Split(filter, "=")
		if len(values) != 2 {
			return nil, fmt.Errorf(httputil.INCORRECT_FORMAT, "filters")
		}
		k, v := values[0], values[1]
		filters[k] = v
	}
	return filters, nil
}
