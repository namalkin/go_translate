package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/namalkin/go_translate/docs"
	"github.com/namalkin/go_translate/pkg/handler"
	"github.com/namalkin/go_translate/pkg/repository"
	"github.com/namalkin/go_translate/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error godotenv.Load(): %s", err.Error())
	}

	// --- MongoDB ---
	db, err := repository.NewMongoDB(repository.MongoConfig{
		URI:      viper.GetString("db.uri"),
		Database: viper.GetString("db.dbname"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize MongoDB: %s", err.Error())
	}
	// --- MongoDB ---

	// --- Redis ---
	redisRepo, err := repository.NewRedisRepo(
		viper.GetString("redis.addr"),
		viper.GetString("redis.password"),
		viper.GetInt("redis.db"))

	if err != nil {
		logrus.Fatalf("failed to initialize Redis: %s", err.Error())
	}
	// --- Redis ---

	repos := repository.NewRepository(db, viper.GetString("db.dbname"), "users", redisRepo)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
