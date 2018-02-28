{{if .Model.IsJoinModel -}}
{{else -}}
{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
var (
	{{$varNameSingular}}Columns               = []string{{"{"}}{{.Model.Columns | columnNames | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
	{{$varNameSingular}}ColumnsWithoutDefault = []string{{"{"}}{{.Model.Columns | filterColumnsByDefault false | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
	{{$varNameSingular}}ColumnsWithDefault    = []string{{"{"}}{{.Model.Columns | filterColumnsByDefault true | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
	{{$varNameSingular}}PrimaryKeyColumns     = []string{{"{"}}{{.Model.PrimaryKey.Columns | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
)

type (
	// {{$modelNameSingular}}Slice is an alias for a slice of pointers to {{$modelNameSingular}}.
	// This should generally be used opposed to []{{$modelNameSingular}}.
	{{$modelNameSingular}}Slice []*{{$modelNameSingular}}
	{{if not .NoHooks -}}
	// {{$modelNameSingular}}Hook is the signature for custom {{$modelNameSingular}} hook methods
	{{$modelNameSingular}}Hook func(context.Context, *{{$modelNameSingular}}) error
	{{- end}}

	{{$varNameSingular}}Query struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	{{$varNameSingular}}Type = reflect.TypeOf(&{{$modelNameSingular}}{})
	{{$varNameSingular}}Mapping = queries.MakeStructMapping({{$varNameSingular}}Type)
	{{$varNameSingular}}PrimaryKeyMapping, _ = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, {{$varNameSingular}}PrimaryKeyColumns)
	{{$varNameSingular}}InsertCacheMut sync.RWMutex
	{{$varNameSingular}}InsertCache = make(map[string]insertCache)
	{{$varNameSingular}}UpdateCacheMut sync.RWMutex
	{{$varNameSingular}}UpdateCache = make(map[string]updateCache)
	{{$varNameSingular}}UpsertCacheMut sync.RWMutex
	{{$varNameSingular}}UpsertCache = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key field that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
{{end -}}
