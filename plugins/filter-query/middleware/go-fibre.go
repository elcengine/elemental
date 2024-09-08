package filter_query_middleware

import (
	"github.com/elcengine/elemental/plugins/filter-query"
	"github.com/gofiber/fiber/v2"
)

func Parse(ctx *fiber.Ctx) error {
	ctx.Locals(CTXKey, filter_query.Parse(string(ctx.Request().URI().QueryString())))
	return ctx.Next()
}
