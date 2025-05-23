package elemental

import (
	"context"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
)

func transaction(fn func(ctx mongo.SessionContext) (any, error), alias *string) (any, error) {
	session, err := GetConnection(lo.FromPtr(alias)).StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.TODO())
	return session.WithTransaction(context.TODO(), fn)
}

// Run a transaction with the given function. This is the simplest form of a transaction which uses the default connection and takes care of session management
func Transaction(fn func(ctx mongo.SessionContext) (any, error)) (any, error) {
	return transaction(fn, nil)
}

// Essentially the same as Transaction, but with an alias which points to the connection to use
func ClientTransaction(alias string, fn func(ctx mongo.SessionContext) (any, error)) (any, error) {
	return transaction(fn, &alias)
}

// TransactionBatch runs a batch of queries in a transaction. If any of the queries fail, the transaction is aborted and all changes are rolled back.
// It returns a slice of results and a slice of errors. The results are in the same order as the queries, and the errors are in the same order as the results.
func TransactionBatch(queries ...ModelInterface[any]) ([]any, []any) {
	sessions := make([]mongo.Session, len(queries))
	var results []any
	var errs []any
	for i, q := range queries {
		session, err := q.Connection().StartSession()
		if err != nil {
			panic(err)
		}
		sessions[i] = session
		session.StartTransaction()
		err = mongo.WithSession(context.Background(), session, func(sessionCtx mongo.SessionContext) error {
			lo.TryCatchWithErrorValue(func() error {
				result := q.Exec(sessionCtx)
				results = append(results, result)
				return nil
			}, func(err any) {
				errs = append(errs, err)
			})
			return nil
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, session := range sessions {
		if len(errs) > 0 {
			session.AbortTransaction(context.Background())
		} else {
			session.CommitTransaction(context.Background())
		}
		session.EndSession(context.Background())
	}
	return results, errs
}
