package enums

type OrganizationType string

const (
	OrgTypeIndividual       OrganizationType = "IE"
	OrgTypeLimitedLiability OrganizationType = "LLC"
	OrgTypeJointStock       OrganizationType = "JSC"
)

type ServiceType string

const (
	ServiceTypeConstruction ServiceType = "Construction"
	ServiceTypeDelivery     ServiceType = "Delivery"
	ServiceTypeManufacture  ServiceType = "Manufacture"
)

var ServiceTypesList = []string{
	string(ServiceTypeDelivery),
	string(ServiceTypeConstruction),
	string(ServiceTypeManufacture),
}

type AuthorType string

const (
	AuthorTypeUser         AuthorType = "User"
	AuthorTypeOrganization AuthorType = "Organization"
)
