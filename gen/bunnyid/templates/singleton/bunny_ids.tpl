import (
	"errors"
    "fmt"
    "strings"

    "github.com/sqlbunny/sqlbunny/runtime/bunnyid"
)

func IDFromString(s string) (bunnyid.ID, error) {
	parts := strings.Split(s, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Wrong parts count, expected 2 got %d", len(parts))
	}
	switch parts[0] {
    {{- range $t := .IDTypes }}
	case "{{$t.Prefix}}":
		return {{$t.Name | titleCase}}FromString(s)
    {{- end}}
	}
	return nil, fmt.Errorf("Unknown ID type %s", parts[0])
}


var idPrefixes = map[string]string{
    {{- range $t := .IDTypes }}
	"{{$t.Prefix}}": "{{$t.Name}}",
    {{- end}}
}
