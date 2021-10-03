package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rabdavinci/loan/data"
)

// Loans is a http.Handler
type Loans struct {
	l *log.Logger
}

// NewLoans creates a loans handler with the given logger
func NewLoans(l *log.Logger) *Loans {
	return &Loans{l}
}

func (l *Loans) GetLoans(rw http.ResponseWriter, r *http.Request) {
	l.l.Println("Handle GET Loans")

	// fetch the loans from the datastore
	ll := data.GetLoans()

	// serialize the list to JSON
	err := ll.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (l *Loans) AddLoan(rw http.ResponseWriter, r *http.Request) {
	l.l.Println("Handle POST Loan")

	loan := r.Context().Value(KeyLoan{}).(data.Loan)
	data.AddLoan(&loan)

	err := loan.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

type KeyLoan struct{}

func (l Loans) MiddlewareValidateLoan(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		loan := data.Loan{}
		err := loan.FromJSON(r.Body)
		if err != nil {
			l.l.Println("[ERROR] deserializing loan", err)
			http.Error(rw, "Error reading loan", http.StatusBadRequest)
			return
		}

		// validate the loan
		err = loan.Validate()
		if err != nil {
			l.l.Println("[ERROR] validating loan", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating loan: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the loan to the context
		ctx := context.WithValue(r.Context(), KeyLoan{}, loan)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
