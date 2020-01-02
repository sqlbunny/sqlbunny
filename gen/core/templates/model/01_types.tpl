{{- $varNameSingular := .Model.Name | singular | camelCase -}}
{{- $modelNameSingular := .Model.Name | singular | titleCase -}}
var (
	{{$varNameSingular}}Columns               = []string{{"{"}}{{modelColumns      .Model | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
	{{$varNameSingular}}PrimaryKeyColumns     = []string{{"{"}}{{modelPKColumns    .Model | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
	{{$varNameSingular}}NonPrimaryKeyColumns  = []string{{"{"}}{{modelNonPKColumns .Model | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
)

type (
	// {{$modelNameSingular}}Slice is an alias for a slice of pointers to {{$modelNameSingular}}.
	// This should generally be used opposed to []{{$modelNameSingular}}.
	{{$modelNameSingular}}Slice []*{{$modelNameSingular}}

	{{$varNameSingular}}Query struct {
		*queries.Query
	}
)

// Cache for insert, update
var (
	{{$varNameSingular}}Type = reflect.TypeOf(&{{$modelNameSingular}}{})
	{{$varNameSingular}}Mapping = queries.MakeStructMapping({{$varNameSingular}}Type)
	{{$varNameSingular}}PrimaryKeyMapping, _ = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, {{$varNameSingular}}PrimaryKeyColumns)
	{{$varNameSingular}}InsertCacheMut sync.RWMutex
	{{$varNameSingular}}InsertCache = make(map[string]insertCache)
	{{$varNameSingular}}UpdateCacheMut sync.RWMutex
	{{$varNameSingular}}UpdateCache = make(map[string]updateCache)
)
