package models

type Endereco struct {
	ID                int    `json: "id"`
	UserID            string `json: "user_id"`
	AddressType       string `json: "address_type"`
	StreetType        string `json: "street_type"`
	StreetName        string `json: "street_name"`
	AddressComplement string `json: "address_complement"`
	Number            string `json: "number"`
	Zipcode           string `json: "zipcode"`
	Neighborhood      string `json: "neighborhood"`
	AddressCity       string `json: "address_city"`
	AddressState      string `json: "address_state"`
	IsActive          bool   `json: "is_active"`
}
