package main

//http://localhost:3000/swagger/index.html
import (
	_ "WebProject_part7/docs" //если не ставить в начале _ то Go определяет этот импорт как неиспользуемый и ругается так что воть
	"WebProject_part7/internal/config"
	"WebProject_part7/internal/repository/mongo"
	"WebProject_part7/internal/service"
	"WebProject_part7/internal/transport/rest/handler"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
	"time"
)

// @title Fiber Swagger Example API
// @version 2.0
// @description This is a sample serever
// @termsOfService http://swagger.io/terms/

// @host localhost:8000
// @BasePath /
// @schemes http
func main() {
	if err := SetupViper(); err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()

	config.SetupSwagger(app)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	mongoDataBase, err := config.SetupMongoDataBase(ctx, cancel)

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := mongo.NewUserRepository(mongoDataBase.Collection("users"))
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	userHandler.InitRoutes(app)

	port := viper.GetString("http.port")
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func SetupViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
