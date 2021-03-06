package refresh

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/clients/feedly"
	"github.com/skatekrak/scribe/clients/vimeo"
	"github.com/skatekrak/scribe/clients/youtube"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/loaders"
	"github.com/skatekrak/scribe/services"
	"github.com/skatekrak/utils/middlewares"
	"gorm.io/gorm"
)

type RefreshQuery struct {
	Types []string `query:"types" validate:"required,dive,eq=vimeo|eq=youtube|eq=rss"`
}

type RefreshSourceQuery struct {
	Force bool `query:"force"`
}

func Route(app *fiber.App, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")
	feedlyCategoryID := os.Getenv("FEEDLY_FETCH_CATEGORY_ID")

	youtubeClient := youtube.New(os.Getenv("YOUTUBE_API_KEY"))
	vimeoClient := vimeo.New(os.Getenv("VIMEO_API_KEY"))
	feedlyClient := feedly.New(os.Getenv("FEEDLY_API_KEY"))
	fetcher := fetchers.New(vimeoClient, youtubeClient, feedlyClient)

	sourceService := services.NewSourceService(db)
	contentService := services.NewContentService(db)
	refreshService := services.NewRefreshService(db, fetcher, feedlyCategoryID)

	controller := &Controller{
		rs:               refreshService,
		ss:               sourceService,
		cs:               contentService,
		fetcher:          fetcher,
		feedlyCategoryID: feedlyCategoryID,
	}
	auth := middlewares.Authorization(apiKey)
	sourceLoader := loaders.SourceLoader(sourceService)

	router := app.Group("refresh")

	router.Post("", auth, middlewares.QueryHandler[RefreshQuery](), controller.RefreshByTypes)
	router.Post("/sync-feedly-sources", auth, controller.RefreshFeedly)
	router.Post("/:sourceID", auth, middlewares.QueryHandler[RefreshSourceQuery](), sourceLoader, controller.RefreshSource)
}
