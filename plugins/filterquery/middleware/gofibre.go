package fqm

import (
	"github.com/elcengine/elemental/plugins/filterquery"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// NewGoFiber is a middleware for Fiber that parses the query string and stores the result in the context.
// It uses the filterquery plugin to parse the query string and apply filters, sorting, lookups, and projections to the final query.
//
// Usage:
//
//	app := fiber.New()
//	app.Use(fqm.NewGoFiber())
//	app.Get("/users", func(ctx *fiber.Ctx) error {
//		q := ctx.Locals(fqm.CtxKey).(fq.Result)
//		users := UserModel.Find(q.Filters).Sort(q.Sorts).Select(q.Select).Populate(q.Include...).ExecTT()
//		return ctx.JSON(users)
//	})
func NewGoFiber(opts ...Options) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		result := fq.Parse(string(ctx.Request().URI().QueryString()))
		if len(opts) > 0 {
			if opts[0].DefaultLimit > 0 {
				result.Limit = lo.CoalesceOrEmpty(result.Limit, opts[0].DefaultLimit)
			}
		}
		ctx.Locals(CtxKey, result)
		return ctx.Next()
	}
}
