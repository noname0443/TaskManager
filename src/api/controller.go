package api

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/noname0443/task_manager/env"
	"github.com/noname0443/task_manager/models"
	"github.com/sethvargo/go-retry"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Controller struct {
	db *gorm.DB
}

func NewController() *Controller {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv(env.POSTGRES_HOST),
		os.Getenv(env.POSTGRES_USER),
		os.Getenv(env.POSTGRES_PASSWORD),
		os.Getenv(env.POSTGRES_DBNAME),
		os.Getenv(env.POSTGRES_PORT),
	)
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

	logrus.Info("migration has started")
	err = db.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("migration has been completed")

	return &Controller{
		db: db,
	}
}
