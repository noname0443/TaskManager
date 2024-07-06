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
	req := getUsersReq{}
	if err := req.fromContext(ctx); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}

	fields := logrus.Fields{
		"limit":   req.limit,
		"offset":  req.offset,
		"filters": req.filters,
	}

	usersDB, err := c.getUsers(req.limit, req.offset, req.filters)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	usersJson, err := usersDBtoJSON(usersDB)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debugln("GetUsers")
	ctx.JSON(http.StatusOK, usersJson)
}

func (c *Controller) getUsers(limit, offset int, filters map[string]string) ([]models.User, error) {
	users := []models.User{}
	if err := c.db.Where(filters).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

type getUsersReq struct {
	limit   int
	offset  int
	filters map[string]string
}

func (req *getUsersReq) fromContext(ctx *gin.Context) (err error) {
	req.limit, err = strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		return fmt.Errorf(httputil.INCORRECT_FORMAT, "limit")
	}

	req.offset, err = strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		return fmt.Errorf(httputil.INCORRECT_FORMAT, "offset")
	}

	req.filters, err = ParseFilters(ctx.Query("filters"))
	if err != nil {
		return err
	}

	return nil
}

type UserJSON struct {
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

func usersDBtoJSON(usersDB []models.User) ([]UserJSON, error) {
	byteArray, err := json.Marshal(usersDB)
	if err != nil {
		return nil, err
	}

	usersJson := []UserJSON{}
	err = json.Unmarshal(byteArray, &usersJson)
	if err != nil {
		return nil, err
	}

	return usersJson, nil
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
