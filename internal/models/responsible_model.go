package models

type ResponsibleDbModel struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id"`
	EmployeeId     string `json:"user_id"`
}

type ResponsibleDtoModel struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organizationId"`
	EmployeeId     string `json:"employeeId"`
}
