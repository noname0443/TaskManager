package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/httputil"
	"github.com/sethvargo/go-retry"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Controller struct {
	db *gorm.DB
}

func NewController() *Controller {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	var db *gorm.DB
	var err error

	logrus.Info("connecting to postgres")

	ctx := context.Background()
	if err := retry.Fibonacci(ctx, 1*time.Second, func(ctx context.Context) error {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logrus.Error(err)
			return retry.RetryableError(err)
		}
		return nil
	}); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("connected to postgres successfully")
	db.AutoMigrate()
	logrus.Info("migration has been completed")

	return &Controller{
		db: db,
	}
}

// CreateUser godoc
//
// @Summary      Creates user
// @Description  Creates user
// @Accept       json
// @Produce      json
// @Param        passportNumber body string true "Create user"
// @Success      200  {int} userId
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/ [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}

type User struct {
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type UpdateUserRequest struct {
	UserId  int  `json:"userId"`
	NewUser User `json:"newUser"`
}

// UpdateUser godoc
//
// @Summary      Updates user
// @Description  Updates user
// @Accept       json
// @Produce      json
// @Param        updateUserRequest body UpdateUserRequest true "updateUserRequest"
// @Success      200  {string}  string	"ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/ [put]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}

// GetUsers godoc
//
// @Summary      Gets users
// @Description  Gets users
// @Accept       json
// @Produce      json
// @Param        offset query int true "Pagination offset"
// @Param        limit query int true "Pagination limit"
// @Param        filters query []string false "Filters"
// @Success      200  {array}  User
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/ [get]
func (c *Controller) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}

// DeleteUser godoc
//
// @Summary      Deletes user
// @Description  Deletes user
// @Accept       json
// @Produce      json
// @Param        userid path int true "User ID"
// @Success      200  {string}  string	"ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /users/{userid} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}

type UserTask struct {
	Task    string `json:"task"`
	Hours   string `json:"hours"`
	Minutes string `json:"minutes"`
	TaskID  int    `json:"taskId"`
}

// GetUserTasks godoc
//
// @Summary      Gets the user's tasks
// @Description  Gets the user's tasks
// @Accept       json
// @Produce      json
// @Param        from query string true "from"
// @Param        to query string true "to"
// @Success      200  {array} UserTask
// @Failure      400  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /users/{userid} [get]
func (c *Controller) GetUserTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}

// UpdateTaskStatus godoc
//
// @Summary      Update the task's status
// @Description  Update the task's status
// @Accept       json
// @Produce      json
// @Param        status query bool true "Status"
// @Param        taskid path int true "TaskID"
// @Success      200  {string} string "ok"
// @Failure      400  {object} httputil.HTTPError
// @Failure      404  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /tasks/{taskid} [put]
func (c *Controller) UpdateTaskStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}

// CreateTask godoc
//
// @Summary      Creates task
// @Description  Creates task
// @Accept       json
// @Produce      json
// @Param        description body string true "Description"
// @Success      200  {int} taskId
// @Failure      400  {object} httputil.HTTPError
// @Failure      404  {object} httputil.HTTPError
// @Failure      500  {object} httputil.HTTPError
// @Router       /tasks/ [post]
func (c *Controller) CreateTask(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, 0)
	if false {
		httputil.NewError(ctx, 400, errors.New("test"))
	}
}
