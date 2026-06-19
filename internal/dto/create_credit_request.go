package dto

type CreateCreditRequest struct {
	AccountID  string  `json:"account_id"`
	Principal  float64 `json:"principal"`
	TermMonths int     `json:"term_months"`
}

func (r *CreateCreditRequest) Validate() error {
	if r.AccountID == "" {
		return NewErrorResponse("Account ID is required")
	}
	if r.Principal <= 0 {
		return NewErrorResponse("Principal must be greater than zero")
	}
	if r.TermMonths <= 0 {
		return NewErrorResponse("Term must be greater than zero months")
	}

	return nil
}
