package schema

type Schema struct {
	Tables map[string]*Table `json:"tables"`
}

func NewSchema() *Schema {
	return &Schema{
		Tables: make(map[string]*Table),
	}
}
