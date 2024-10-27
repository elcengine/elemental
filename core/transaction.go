package elemental

import (
	"context"
	"github.com/elcengine/elemental/connection"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
)

func transaction(fn func(ctx mongo.SessionContext) (interface{}, error), alias *string) (interface{}, error) {
	session, err := lo.ToPtr(e_connection.GetConnection(lo.FromPtr(alias))).StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.TODO())
	return session.WithTransaction(context.TODO(), fn)
}

// Run a transaction with the given function. This is the simplest form of a transaction which uses the default connection and takes care of session management
func Transaction(fn func(ctx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	return transaction(fn, nil)
}

// Essentially the same as Transaction, but with an alias which points to the connection to use
func ClientTransaction(alias string, fn func(ctx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	return transaction(fn, &alias)
}

func TransactionBatch(queries ...Model[any]) []interface{} {	var sessions []mongo.Session
	var results []interface{}
	var errCount int
	for _, query := range queries {
		 func(q Model[any]) {
			connection, _ := lo.Coalesce(lo.FromPtr(q.temporaryConnection), q.Schema.Options.Connection)
			session, err := lo.ToPtr(e_connection.GetConnection(connection)).StartSession()
			if err != nil {
				panic(err)
			}
			sessions = append(sessions, session)
			session.StartTransaction()
			err = mongo.WithSession(context.Background(), session, func(sessionCtx mongo.SessionContext) error {
				lo.TryCatch(func() error {
					result := q.Exec(sessionCtx)
					results = append(results, result)
					return nil
				}, func() {
					errCount++
				})
				return nil
			})
			if err != nil {
				errCount++
			}
		}(query)
	}
	for _, session := range sessions {
		if errCount > 0 {
			session.AbortTransaction(context.Background())
		} else {
			session.CommitTransaction(context.Background())
		}
		session.EndSession(context.Background())
	}
	return results
}
