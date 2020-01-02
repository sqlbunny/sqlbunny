package schema

type Schema struct {
	Tables map[string]*Table `json:"tables"`
}

func New() *Schema {
	return &Schema{
		Tables: make(map[string]*Table),
	}
}
