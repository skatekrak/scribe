package refresh

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/loaders"
	"github.com/skatekrak/scribe/services"
	"github.com/skatekrak/utils/middlewares"
)

type Controller struct {
	rs               *services.RefreshService
	ss               *services.SourceService
	cs               *services.ContentService
	fetcher          *fetchers.Fetcher
	feedlyCategoryID string
}

// Refresh sources by their types
// @Summary   Refresh sources by there types
// @Security  ApiKeyAuth
// @Tags      refresh
// @Success   200    {array}   []model.Content
// @Failure   500    {object}  api.JSONError
// @Param     types  query     []string  true  "Type of sources to refresh"  Enums(rss,vimeo,youtube)
// @Router    /refresh [patch]
func (c *Controller) RefreshByTypes(ctx *fiber.Ctx) error {
	query := ctx.Locals(middlewares.QUERY).(RefreshQuery)

	contents, err := c.rs.RefreshByTypes(query.Types)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(contents)
}

// Refresh a given sources
// @Summary   Refresh a given source
// @Security  ApiKeyAuth
// @Tags      refresh
// @Success   200       {array}   []model.Source
// @Failure   500       {object}  api.JSONError
// @Param     sourceID  path      string  true   "Source ID"
// @Param     force     query     bool    false  "Will override content attributes"
// @Router    /refresh/{sourceID} [patch]
func (c *Controller) RefreshSource(ctx *fiber.Ctx) error {
	source := loaders.GetSource(ctx)

	query := ctx.Locals(middlewares.QUERY).(RefreshSourceQuery)

	contents, errs := c.rs.RefreshBySource(source, query.Force)
	if errs != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errs)
	}

	return ctx.Status(fiber.StatusOK).JSON(contents)
}

// Refresh feedly sources
// @Summary   Query sources used in feedly and add missing ones in Scribe
// @Security  ApiKeyAuth
// @Tags      refresh
// @Success   200  {array}   []model.Source
// @Failure   500  {object}  api.JSONError
// @Router    /refresh/sync-feedly [patch]
func (c *Controller) RefreshFeedly(ctx *fiber.Ctx) error {
	sources, err := c.rs.RefreshFeedlySource()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(sources)
}
