package helpers

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"bcpayslip/models"

	"github.com/jung-kurt/gofpdf"
)

// ImageToBase64 Convert url image to base64 encoding ...
func ImageToBase64(url string) string {
	res, _ := http.Get(url)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	imgBase64Str := base64.StdEncoding.EncodeToString(bodyBytes)
	return imgBase64Str
}

// ConvertFormDate Converts html date strings to a date type format and returns it ...
func ConvertFormDate(value string) reflect.Value {
	s, _ := time.Parse("2006-01-02", value)
	return reflect.ValueOf(s.UTC())
}

// GeneratePayslipPDF generate PDF for payslip ...
func GeneratePayslipPDF(payslip *models.Payslip) error {
	salary := payslip.GrossAnnualSalary
	var floatInt int
	floatInt = 2
	var floatType int
	floatType = 64
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetX(-60)
	pdf.SetFont("Arial", "", 16)
	pdf.SetTextColor(26, 162, 251)
	pdf.Cell(30, 0, "BEAUTIFUL ")
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(30, 0, " CODE")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(10)
	pdf.SetXY(100, 30)
	pdf.Line(10, 20, 200, 20)
	pdf.Line(10, 40, 200, 40)
	pdf.Cell(100, 0, "Pay Slip")
	pdf.Line(10, 20, 10, 40)
	pdf.Line(200, 20, 200, 40)
	pdf.SetXY(100, 40)
	pdf.Line(10, 40, 200, 40)
	pdf.Line(10, 70, 200, 70)
	pdf.Line(10, 40, 10, 70)
	pdf.Line(200, 40, 200, 70)
	pdf.SetXY(20, 40)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Pay Period: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.Month.Format(time.RFC1123)[8:16])
	pdf.SetXY(100, 40)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Pay Date: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.Day.Format(time.RFC1123)[5:16])
	pdf.SetXY(20, 50)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Employee Name: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.Name)
	pdf.SetXY(100, 50)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Position: ")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, payslip.Position)
	pdf.SetXY(20, 60)
	if strings.Replace(payslip.EmployeeNo, " ", "", -1) != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(40, 10, "Employee No: ")
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(40, 10, payslip.EmployeeNo)
	}
	pdf.SetXY(100, 70)
	pdf.Line(10, 70, 200, 70)
	pdf.Line(10, 120, 200, 120)
	pdf.Line(10, 70, 10, 120)
	pdf.Line(120, 70, 120, 120)
	pdf.Line(200, 70, 200, 120)
	pdf.SetXY(20, 70)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(70, 10, "Earnings & Allowances")
	pdf.Cell(30, 10, "INR")
	pdf.SetXY(20, 80)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(70, 10, "Basic Salary")
	pdf.Cell(30, 10, strconv.FormatFloat((salary*0.6), 'f', floatInt, floatType))
	pdf.SetXY(20, 90)
	pdf.Cell(70, 10, "House Rent Allowance")
	pdf.Cell(30, 10, strconv.FormatFloat((salary*0.2), 'f', floatInt, floatType))
	pdf.SetXY(20, 100)
	pdf.Cell(70, 10, "Spcial / Conv Allowance")
	pdf.Cell(30, 10, strconv.FormatFloat((salary*0.15), 'f', floatInt, floatType))
	pdf.SetXY(20, 110)
	pdf.Cell(70, 10, "Other Allowance")
	pdf.Cell(30, 10, strconv.FormatFloat((salary*0.05), 'f', floatInt, floatType))
	pdf.SetXY(120, 70)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Deductions")
	pdf.Cell(20, 10, "INR")
	pdf.SetXY(120, 80)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Income Tax")
	pdf.Cell(20, 10, strconv.FormatFloat(salary-payslip.AmountReceivedBank, 'f', floatInt, floatType))
	pdf.SetXY(120, 90)
	pdf.Cell(40, 10, "Advance")
	pdf.Cell(20, 10, "0.00")
	pdf.SetXY(120, 100)
	pdf.Cell(40, 10, "Profession Tax")
	pdf.Cell(20, 10, "0.00")
	pdf.SetXY(100, 120)
	pdf.Line(10, 120, 200, 120)
	pdf.Line(10, 160, 200, 160)
	pdf.Line(10, 120, 10, 160)
	pdf.Line(120, 120, 120, 160)
	pdf.Line(200, 120, 200, 160)
	pdf.SetXY(20, 120)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(20, 10, "Bank Account: ")
	pdf.SetXY(20, 130)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Account No: ")
	pdf.Cell(50, 10, payslip.AccountNo)
	pdf.SetXY(20, 140)
	pdf.Cell(40, 10, "IFSC Code: ")
	pdf.Cell(50, 10, payslip.IFSCCode)
	pdf.SetXY(120, 120)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Pay Summary")
	pdf.Cell(30, 10, "INR")
	pdf.SetXY(120, 130)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Total Gross")
	pdf.Cell(20, 10, strconv.FormatFloat((salary), 'f', floatInt, floatType))
	pdf.SetXY(120, 140)
	pdf.Cell(40, 10, "Deductions")
	pdf.Cell(20, 10, strconv.FormatFloat(((salary)-payslip.AmountReceivedBank), 'f', floatInt, floatType))
	pdf.SetXY(120, 150)
	pdf.Cell(40, 10, "NET PAY")
	pdf.Cell(20, 10, strconv.FormatFloat(payslip.AmountReceivedBank, 'f', floatInt, floatType))
	pdf.SetXY(10, 160)
	pdf.Line(10, 190, 200, 190)
	pdf.Line(10, 160, 10, 190)
	pdf.Line(200, 160, 200, 190)
	pdf.SetXY(75, 170)
	pdf.Cell(150, 10, "(*) denotes back pay adjustment")
	pdf.SetXY(75, 180)
	pdf.Cell(150, 10, "Computer Generated Form does not require signature")
	err := pdf.OutputFileAndClose("media/" + payslip.PayslipID + ".pdf")
	return err
}
