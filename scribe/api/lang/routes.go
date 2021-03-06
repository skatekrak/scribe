package lang

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/loaders"
	"github.com/skatekrak/scribe/services"
	"github.com/skatekrak/utils/middlewares"
	"gorm.io/gorm"
)

type CreateBody struct {
	IsoCode  string `json:"isoCode" binding:"required,len=2"`
	ImageURL string `json:"imageURL" binding:"required"`
}

type UpdateBody struct {
	ImageURL string `json:"imageURL" binding:"required"`
}

type LangUri struct {
	IsoCode string `uri:"isoCode" binding:"required"`
}

func Route(app *fiber.App, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")

	langService := services.NewLangService(db)
	controller := &Controller{
		s: langService,
	}

	auth := middlewares.Authorization(apiKey)
	langLoader := loaders.LangLoader(langService)

	router := app.Group("/langs")
	router.Get("", controller.FindAll)

	router.Post("", auth, middlewares.JSONHandler[CreateBody](), controller.Create)
	router.Patch("/:isoCode", auth, langLoader, middlewares.JSONHandler[UpdateBody](), controller.Update)
	router.Delete("/:isoCode", auth, langLoader, controller.Delete)
}
