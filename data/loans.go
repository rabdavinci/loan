package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator"
)

type Loan struct {
	ID        int     `json:"id" validate:"unique"`
	Product   string  `json:"product" validate:"product"`
	Phone     string  `json:"phone" validate:"e164"`
	Month     int     `json:"month" validate:"range"`
	Price     float32 `json:"price" validate:"gt=0"`
	CreatedOn string  `json:"-"`
	UpdatedOn string  `json:"-"`
}

type Product struct {
	Name  string
	Range []int
}

func (l *Loan) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("product", validateProduct)
	validate.RegisterValidation("range", validateRange)
	return validate.Struct(l)
}

func validateProduct(fl validator.FieldLevel) bool {
	for _, v := range productList {
		if fl.Field().String() == v.Name {
			return true
		}
	}

	return false
}

func validateRange(fl validator.FieldLevel) bool {
	for _, v := range rangeList {
		if fl.Field().Int() == int64(v) {
			return true
		}
	}

	return false
}

type Loans []*Loan

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
func (l *Loans) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(l)
}

func (l *Loan) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(l)
}

func GetLoans() Loans {
	return loansList
}

func AddLoan(l *Loan) {
	l.ID = getNextID()
	loansList = append(loansList, l)
}

var ErrLoanNotFound = fmt.Errorf("Loan not found")

func getNextID() int {
	ll := loansList[len(loansList)-1]
	return ll.ID + 1
}

// loansList is a hard coded list of loans for this
// example data source
var rangeList = []int{3, 6, 9, 12, 18, 24}
var productList = []*Product{
	&Product{
		Name:  "Смартфон",
		Range: []int{3, 6, 9},
	},
	&Product{
		Name:  "Компьютер",
		Range: []int{3, 6, 9, 12},
	},
	&Product{
		Name:  "Телевизор",
		Range: []int{3, 6, 9, 12, 18},
	},
}
var loansList = []*Loan{
	&Loan{
		ID:        1,
		Product:   "Смартфон",
		Phone:     "+998995881375",
		Month:     3,
		Price:     1000,
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
	&Loan{
		ID:        2,
		Product:   "Смартфон",
		Phone:     "+998995881375",
		Month:     24,
		Price:     1000,
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
}
