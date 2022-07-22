package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/skatekrak/carrelage/auth"
	"github.com/skatekrak/carrelage/models"
	"github.com/skatekrak/utils/database"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	db, err := database.Open(os.Getenv("POSTGRESQL_ADDON_URI"))
	if err != nil {
		log.Fatalf("unable to open database: %s", err)
	}

	if err = db.AutoMigrate(&models.User{}, &models.UserSubscription{}, &models.Profile{}); err != nil {
		log.Fatalf("unable to migrate database: %s", err)
	}

	auth.InitSuperTokens()

	app := fiber.New()

	// Setup middlewares
	app.Use(logger.New())

	corsHeaders := []string{"content-type"}
	corsHeaders = append(corsHeaders, supertokens.GetAllCORSHeaders()...)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ORIGINS"),
		AllowHeaders:     strings.Join(corsHeaders, ","),
		AllowCredentials: true,
	}))

	app.Use(compress.New())
	app.Use(recover.New())

	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))

	// Start server
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalln("Error listening")
		log.Fatalln(err)
	}
}