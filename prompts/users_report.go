package qen

import (
	"bytes"
	"fmt"
	"text/template"
)

type UsersReportInput struct {
	HasPremiumUsers bool
	PremiumUsers []UsersReportInputPremiumUsers
	GeneratedBy string
	GenerationDate string
	Users []UsersReportInputUsers
	TotalUsers string
	ActiveUsers string
	InactiveUsers string
}

type UsersReportInputUsers struct {
	ID string
	FirstName string
	LastName string
	Email string
	Role string
	IsActive bool
}

type UsersReportInputPremiumUsers struct {
	FirstName string
	LastName string
	MemberSince string
}

func UsersReport(input UsersReportInput) (string, error) {
	tmpl, err := template.ParseFiles("prompts/users_report.gohtml")
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buf.String(), nil
}
