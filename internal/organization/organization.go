package organization

type Organization struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Address            string `json:"address"`
	ContactInformation string `json:"contact_information"`
}

// New creates a new instance of Organization with the given details.
func New(name string, address string, contactInformation string) Organization {
	return Organization{
		Name:               name,
		Address:            address,
		ContactInformation: contactInformation,
	}
}
