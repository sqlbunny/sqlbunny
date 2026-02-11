package schema

type Database struct {
	Schemas map[string]*Schema `json:"schemas"`
}

func NewDatabase() *Database {
	return &Database{
		Schemas: make(map[string]*Schema),
	}
}
