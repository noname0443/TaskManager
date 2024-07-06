package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/sirupsen/logrus"
)

// GetUserTasks godoc
//
// @Summary      Gets the user's tasks
// @Description  Gets the user's tasks
// @Accept       json
// @Produce      json
// @Param        to query string false "YYYY-MM-DDThh:mm:ss.SSSZ rfc3339"
// @Param        from query string false "YYYY-MM-DDThh:mm:ss.SSSZ rfc3339"
// @Param        offset query int true "Pagination offset"
// @Param        limit query int true "Pagination limit"
// @Param        userId path int true "UserID"
// @Success      200  {array} UserTask
// @Failure      400  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /users/{userId} [get]
func (c *Controller) GetUserTasks(ctx *gin.Context) {
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

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, fmt.Errorf(httputil.INCORRECT_FORMAT, "userId"))
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", ctx.Query("from"))
	if err != nil {
		logrus.Debug(err)
		from = time.Time{}
	}
	from = from.Round(time.Microsecond)

	to, err := time.Parse("2006-01-02T15:04:05.000Z", ctx.Query("to"))
	if err != nil {
		logrus.Debug(err)
		to = time.Now()
	}
	to = to.Round(time.Microsecond)

	if to.Before(from) {
		logrus.Debug("'to' is before 'from'")
		ctx.JSON(200, []TaskJSON{})
		return
	}

	fields := logrus.Fields{
		"from":   from,
		"to":     to,
		"limit":  limit,
		"offset": offset,
		"userId": userId,
	}

	tasksDB := []TaskDB{}
	rows, err := c.db.Raw(`SELECT id, "userId", description, status, coalesce(EXTRACT(epoch FROM SUM(estimated_time) * 1000000000)::BIGINT, 0) as estimated_time, start FROM tasks A LEFT JOIN (
		SELECT "taskId", (LEAST(end_interval, ?) - GREATEST(begin_interval, ?)) as estimated_time FROM time_spents WHERE begin_interval < ? AND end_interval > ?
	) B ON A.id = B."taskId" WHERE "userId" = ? GROUP BY id, created_at, updated_at, deleted_at, "userId", description, start, status ORDER BY estimated_time DESC, id ASC LIMIT ? OFFSET ?;`, to, from, to, from, userId, limit, offset).Rows()
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}
	for rows.Next() {
		taskDb := TaskDB{}
		c.db.ScanRows(rows, &taskDb)
		tasksDB = append(tasksDB, taskDb)
	}

	byteArray, err := json.Marshal(tasksDB)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	tasksJson := []TaskJSON{}
	err = json.Unmarshal(byteArray, &tasksJson)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	for i := 0; i < len(tasksJson); i++ {
		if tasksDB[i].Status && to.After(tasksDB[i].Start) {
			tasksDB[i].EstimatedTime += to.Sub(tasksDB[i].Start)
		}
		tasksJson[i].Seconds = uint(tasksDB[i].EstimatedTime.Seconds())
		tasksJson[i].Minutes = uint(tasksDB[i].EstimatedTime.Minutes())
		tasksJson[i].Hours = uint(tasksDB[i].EstimatedTime.Hours())
	}

	logrus.WithFields(fields).Debug("GetUserTasks")
	ctx.JSON(200, tasksJson)
}

type TaskDB struct {
	ID            uint          `json:"taskId" gorm:"column:id"`
	UserID        uint          `json:"userId" gorm:"column:userId"`
	Description   string        `json:"description" gorm:"column:description"`
	EstimatedTime time.Duration `json:"estimated_time" gorm:"column:estimated_time"`
	Start         time.Time     `json:"start" gorm:"column:start"`
	Status        bool          `json:"status" gorm:"column:status"`
}

type TaskJSON struct {
	UserID      uint   `json:"userId"`
	ID          uint   `json:"taskId"`
	Description string `json:"description"`
	Hours       uint   `json:"hours"`
	Minutes     uint   `json:"minutes"`
	Seconds     uint   `json:"seconds"`
	Status      bool   `json:"status"`
}
