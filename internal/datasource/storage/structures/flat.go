package structures

type Flat struct {
	Id      int    `db:"id" json:"id,omitempty"`
	HouseId int    `db:"house_id" json:"house_id,omitempty"`
	Price   int    `db:"price" json:"price,omitempty"`
	Rooms   int    `db:"rooms" json:"rooms,omitempty"`
	Status  string `db:"status" json:"status,omitempty"`
}
