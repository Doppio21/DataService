package schema

type PutRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type GetRequest struct {
	ID      int
	Name    string
	Surname string
	Age     int
	Gender  string
	Country string
	Count   int
	Offset  int
}

type GetResponse struct {
	Persons []PersonInfo
}

type PersonInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Gender  string `json:"gender"`
}
