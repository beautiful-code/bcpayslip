package models

import "time"

type (
	// User type represents the registered user. ...
	User struct {
		UserID      string `json:"userid"`
		FirstName   string `json:"firstname"`
		LastName    string `json:"lastname"`
		Email       string `json:"email"`
		AccessToken string `json:"token,omitempty"`
		Avatar      string `json:"avatar"`
	}
	// Message Flash message Struct ...
	Message struct {
		Value string
	}
	// Kwargs Pass keyword args to urls ...
	Kwargs struct {
		Key   string
		Value string
	}
	// Payslip ...
	Payslip struct {
		PayslipID          string    `json:"id"`
		Name               string    `json:"name"`
		Requestor          User      `bson:"requestor" json:"requestor"`
		Approver           User      `bson:"approver" json:"approver"`
		RequestedOn        time.Time `json:requestedon`
		Day                time.Time `json:"day"`
		Month              time.Time `json:"month"`
		GrossAnnualSalary  float64   `json:"salary"`
		AmountReceivedBank float64   `json:"amount"`
		TDS                float64   `json:"tds"`
		AccountNo          string    `json:"accountno"`
		IFSCCode           string    `json:"ifsccode"`
		Position           string    `json:"position"`
		EmployeeNo         string    `json:"employeeno"`
		Status             int       `json:"status"`
	}
)
