package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	_ "time/tzdata"

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
// @Success      200  {array} TaskJSON
// @Failure      400  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /users/{userId} [get]
func (c *Controller) GetUserTasks(ctx *gin.Context) {
	req := getUserTasksReq{}
	if err := req.fromContext(ctx); err != nil {
		logrus.Debug(err)
		httputil.NewError(ctx, 400, err)
		return
	}
	if req.to.Before(req.from) {
		logrus.Debug("'to' is before 'from'")
		ctx.JSON(200, []TaskJSON{})
		return
	}

	fields := logrus.Fields{
		"from":   req.from,
		"to":     req.to,
		"limit":  req.limit,
		"offset": req.offset,
		"userId": req.userId,
	}

	tasksDB, err := c.getTimeSpents(req.userId, req.limit, req.offset, req.from, req.to)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	tasksJson, err := tasksDBtoJSON(tasksDB, req.to)
	if err != nil {
		logrus.WithFields(fields).Warn(err)
		httputil.NewError(ctx, 500, fmt.Errorf(httputil.SOMETHING_WENT_WRONG))
		return
	}

	logrus.WithFields(fields).Debug("GetUserTasks")
	ctx.JSON(200, tasksJson)
}

type getUserTasksReq struct {
	limit  int
	offset int
	from   time.Time
	to     time.Time
	userId int
}

func (req *getUserTasksReq) fromContext(ctx *gin.Context) (err error) {
	req.userId, err = strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		return fmt.Errorf(httputil.INCORRECT_FORMAT, "userId")
	}

	req.limit, err = strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		return fmt.Errorf(httputil.INCORRECT_FORMAT, "limit")
	}

	req.offset, err = strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		return fmt.Errorf(httputil.INCORRECT_FORMAT, "offset")
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	req.from, err = time.ParseInLocation(rfc3339, ctx.Query("from"), loc)
	if err != nil {
		logrus.Debug(err)
		req.from = time.Time{}
	}
	req.from = req.from.Round(time.Microsecond)

	req.to, err = time.ParseInLocation(rfc3339, ctx.Query("to"), loc)
	if err != nil {
		logrus.Debug(err)
		req.to = time.Now()
	}
	req.to = req.to.Round(time.Microsecond)

	return nil
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

func (c *Controller) getTimeSpents(userId, limit, offset int, from, to time.Time) ([]TaskDB, error) {
	tasksDB := []TaskDB{}
	rows, err := c.db.Raw(sumTimeSpentIntervals, to, from, to, from, userId, limit, offset).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		taskDb := TaskDB{}
		err := c.db.ScanRows(rows, &taskDb)
		if err != nil {
			return nil, err
		}

		if taskDb.Status && taskDb.Start.Before(to.Add(-time.Second)) ||
			taskDb.EstimatedTime > time.Second {
			tasksDB = append(tasksDB, taskDb)
		}
	}
	return tasksDB, nil
}

func tasksDBtoJSON(tasksDB []TaskDB, to time.Time) ([]TaskJSON, error) {
	byteArray, err := json.Marshal(tasksDB)
	if err != nil {
		return nil, err
	}

	tasksJson := []TaskJSON{}
	err = json.Unmarshal(byteArray, &tasksJson)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tasksJson); i++ {
		if tasksDB[i].Status && to.After(tasksDB[i].Start) {
			tasksDB[i].EstimatedTime += to.Sub(tasksDB[i].Start)
		}
		tasksJson[i].Seconds = uint(tasksDB[i].EstimatedTime.Seconds()) % 60
		tasksJson[i].Minutes = uint(tasksDB[i].EstimatedTime.Minutes()) % 60
		tasksJson[i].Hours = uint(tasksDB[i].EstimatedTime.Hours()) % 60
	}
	return tasksJson, nil
}

const (
	rfc3339               = "2006-01-02T15:04:05.000Z"
	sumTimeSpentIntervals = `
	SELECT id, "userId", description, status, coalesce(EXTRACT(epoch FROM SUM(estimated_time) * 1000000000)::BIGINT, 0) as estimated_time, start
	FROM tasks A LEFT JOIN (
		SELECT "taskId", (LEAST(end_interval, ?) - GREATEST(begin_interval, ?)) as estimated_time
		FROM time_spents
		WHERE begin_interval <= ? AND end_interval >= ?
	) B ON A.id = B."taskId"
	WHERE "userId" = ? GROUP BY id, created_at, updated_at, deleted_at, "userId", description, start, status
	ORDER BY estimated_time DESC, id ASC LIMIT ? OFFSET ?;`
)
