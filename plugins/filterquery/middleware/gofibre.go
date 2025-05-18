package fqm

import (
	"github.com/elcengine/elemental/plugins/filterquery"
	"github.com/gofiber/fiber/v2"
)

// NewGoFiber is a middleware for Fiber that parses the query string and stores the result in the context.
// It uses the filterquery plugin to parse the query string and apply filters, sorting, lookups, and projections to the final query.
//
// Usage:
//
//	app := fiber.New()
//	app.Use(fqm.NewGoFiber())
//	app.Get("/users", func(ctx *fiber.Ctx) error {
//		q := ctx.Locals(fqm.CtxKey).(fq.FilterQueryResult)
//		users := UserModel.Find(q.Filters).Sort(q.Sorts).Select(q.Select).Populate(q.Include...).ExecTT()
//		return ctx.JSON(users)
//	})
func NewGoFiber() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(CtxKey, fq.Parse(string(ctx.Request().URI().QueryString())))
		return ctx.Next()
	}
}
