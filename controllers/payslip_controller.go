package controllers

import (
	"net/http"
	"time"

	"bcpayslip/helpers"
	"bcpayslip/models"
	"bcpayslip/store"
	"bcpayslip/templates"
	"bcpayslip/urls"
	"bcpayslip/utils"

	"github.com/gorilla/context"
	"github.com/gorilla/schema"
)

// PayslipController ...
func PayslipController(res http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	controllerTemplate := templates.PayslipTemplate
	if req.Method == "GET" {
		utils.CustomTemplateExecute(res, req, controllerTemplate, data)
	}
	if req.Method == "POST" {
		err := req.ParseForm()
		payslip := new(models.Payslip)
		decoder := schema.NewDecoder()
		decoder.RegisterConverter(time.Time{}, helpers.ConvertFormDate)
		err = decoder.Decode(payslip, req.Form)
		if err != nil {
			http.Redirect(res, req, urls.HomePath, http.StatusSeeOther)
		}
		user, _ := store.GetUser(context.Get(req, "userid").(string))
		payslip.Requestor = user
		payslip.RequestedOn = time.Now()
		payslip.Status = 0
		payslip.PayslipID = user.UserID
		helpers.GeneratePayslipPDF(payslip)
		http.Redirect(res, req, "/media/"+payslip.PayslipID+".pdf", http.StatusSeeOther)
	}
}
