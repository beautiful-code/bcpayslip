package main

import (
	"testing"

	"bcpayslip/helpers"
	"bcpayslip/models"
)

func TestPDF(t *testing.T) {
	payslip := new(models.Payslip)
	payslip.PayslipID = "123456789012"
	payslip.GrossAnnualSalary = 660000
	err := helpers.GeneratePayslipPDF(payslip)
	if err != nil {
		t.Errorf("PDF error: %s", err)
	}
}
