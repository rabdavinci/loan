package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator"
)

type Loan struct {
	ID         int     `json:"id"`
	Product    string  `json:"product" validate:"product"`
	Phone      string  `json:"phone" validate:"e164"`
	Month      int     `json:"month" validate:"range"`
	Price      float32 `json:"price" validate:"gt=0"`
	TotalPrice float32
	CreatedOn  string `json:"-"`
	UpdatedOn  string `json:"-"`
}

type Product struct {
	Name       string
	PeriodFree int
	Percent    int
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

var rangeList = []int{3, 6, 9, 12, 18, 24}

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

func (l *Loan) ToJSON(w io.Writer) error {
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
	l.TotalPrice = getTotalPrice(l.Product, l.Month, l.Price)
	l.CreatedOn = time.Now().UTC().String()
	l.UpdatedOn = time.Now().UTC().String()
	loansList = append(loansList, l)
}

func getTotalPrice(product string, month int, price float32) float32 {
	prod, _ := findProduct(product)

	if month <= prod.PeriodFree {
		return price
	}

	return price + price*float32((month-prod.PeriodFree)*prod.Percent)/(100*3)
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(name string) (*Product, error) {
	for _, p := range productList {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, ErrProductNotFound
}

var ErrLoanNotFound = fmt.Errorf("Loan not found")

func getNextID() int {
	ll := loansList[len(loansList)-1]
	return ll.ID + 1
}

// loansList is a hard coded list of loans for this
// example data source

var productList = []*Product{
	&Product{
		Name:       "Смартфон",
		PeriodFree: 9,
		Percent:    3,
	},
	&Product{
		Name:       "Компьютер",
		PeriodFree: 12,
		Percent:    4,
	},
	&Product{
		Name:       "Телевизор",
		PeriodFree: 18,
		Percent:    5,
	},
}

// hardcoding example loans
var loansList = []*Loan{
	&Loan{
		ID:         1,
		Product:    "Смартфон",
		Phone:      "+998995881375",
		Month:      3,
		Price:      1000,
		TotalPrice: getTotalPrice("Смартфон", 3, 1000),
		CreatedOn:  time.Now().UTC().String(),
		UpdatedOn:  time.Now().UTC().String(),
	},
	&Loan{
		ID:         2,
		Product:    "Смартфон",
		Phone:      "+998995881375",
		Month:      24,
		Price:      1000,
		TotalPrice: getTotalPrice("Смартфон", 24, 1000),
		CreatedOn:  time.Now().UTC().String(),
		UpdatedOn:  time.Now().UTC().String(),
	},
}
