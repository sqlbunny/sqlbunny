package def

import "github.com/kernelpayments/sqlbunny/schema"

type Config struct {
	Items  []ConfigItem
	Schema *schema.Schema
}
