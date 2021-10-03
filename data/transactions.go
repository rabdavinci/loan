package data

// здесь описал структуру погашения кредита
type Transaction struct {
	ID        int     `json:"id"`
	LoanID    int     `json:"loanID"`
	Amount    float32 `json:"amount" validate:"gt=0"`
	CreatedOn string  `json:"-"`
	UpdatedOn string  `json:"-"`
}
