package main

import (
	"fmt"
	"os"

	_ "github.com/noname0443/task_manager/docs"
	"github.com/noname0443/task_manager/env"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/noname0443/task_manager/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	r := gin.Default()

	LoadLoggerLevel()

	c := api.NewController()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users/:userId", c.GetUserTasks)
		v1.GET("/users", c.GetUsers)

		v1.POST("/tasks", c.CreateTask)
		v1.POST("/users", c.CreateUser)

		v1.PUT("/tasks/:taskId", c.UpdateTaskStatus)
		v1.PUT("/users/:userId", c.UpdateUser)

		v1.DELETE("/users/:userId", c.DeleteUser)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(fmt.Sprintf(":%s", os.Getenv(env.SERVICE_PORT)))
}

func LoadLoggerLevel() {
	log_level := os.Getenv(env.LOG)

	if log_level == env.LOG_DEBUG {
		logrus.SetLevel(logrus.DebugLevel)
	} else if log_level == env.LOG_INFO {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.Fatal("unknown log level")
	}
}
