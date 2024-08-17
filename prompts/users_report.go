package typlmpts

import (
	"bytes"
	"fmt"
	"text/template"
)

type UsersReportInput struct {
	Users           []UsersReportInputUsers
	TotalUsers      string
	ActiveUsers     string
	InactiveUsers   string
	HasPremiumUsers bool
	PremiumUsers    []UsersReportInputPremiumUsers
	GeneratedBy     string
	GenerationDate  string
}

type UsersReportInputUsers struct {
	LastName  string
	Email     string
	Role      string
	IsActive  bool
	ID        string
	FirstName string
}

type UsersReportInputPremiumUsers struct {
	LastName    string
	MemberSince string
	FirstName   string
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
