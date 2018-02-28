import (
	"bytes"
	"context"
	"database/sql"
    "encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/KernelPay/sqlboiler/boil"
	"github.com/KernelPay/sqlboiler/boil/queries"
	"github.com/KernelPay/sqlboiler/boil/qm"
	"github.com/KernelPay/sqlboiler/boil/strmangle"
	"github.com/pkg/errors"
)
{{- range $import := .Imports }}
import {{$import}}
{{- end -}}

{{- define "relationship_to_one_struct_helper" -}}
{{- end -}}

{{- $dot := . -}}
{{- $modelNameSingular := .Model.Name | singular -}}
{{- $modelName := $modelNameSingular | titleCase -}}
{{- $modelNameCamel := $modelNameSingular | camelCase -}}

// {{$modelName}} is an object representing the database model.
type {{$modelName}} struct {
	{{range $field := .Model.Fields }}
	{{- if eq $dot.StructTagCasing "camel"}}
	{{titleCase $field.Name}} {{$field.TypeGo}} `{{generateTags $dot.Tags $field.Name}}boil:"{{$field.Name}}{{ if $field.IsStruct }},{{end}}" json:"{{$field.Name | camelCase}}{{if $field.Nullable}},omitempty{{end}}" toml:"{{$field.Name | camelCase}}" yaml:"{{$field.Name | camelCase}}{{if $field.Nullable}},omitempty{{end}} {{$field.Tag}}`
	{{- else -}}
	{{titleCase $field.Name}} {{$field.TypeGo}} `{{generateTags $dot.Tags $field.Name}}boil:"{{$field.Name}}{{ if $field.IsStruct }},{{end}}" json:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}" toml:"{{$field.Name}}" yaml:"{{$field.Name}}{{if $field.Nullable}},omitempty{{end}}" {{$field.Tag}}`
	{{end -}}
	{{end -}}
	{{- if .Model.IsJoinModel -}}
	{{- else}}
	R *{{$modelNameCamel}}R `{{generateIgnoreTags $dot.Tags}}boil:"-" json:"-" toml:"-" yaml:"-"`
	L {{$modelNameCamel}}L `{{generateIgnoreTags $dot.Tags}}boil:"-" json:"-" toml:"-" yaml:"-"`
	{{end -}}
}

var {{$modelName}}Columns = struct {
	{{range $column := .Model.Columns -}}
	{{titleCase $column.Name}} string
	{{end -}}
}{
	{{range $column := .Model.Columns -}}
	{{titleCase $column.Name | }}: "{{$column.Name}}",
	{{end -}}
}

{{- if .Model.IsJoinModel -}}
{{- else}}
// {{$modelNameCamel}}R is where relationships are stored.
type {{$modelNameCamel}}R struct {
	{{range .Model.ForeignKeys -}}
	{{- $txt := txtsFromFKey $dot.Models $dot.Model . -}}
	{{$txt.Function.NameGo}} *{{$txt.ForeignModel.NameGo}}
	{{end -}}

	{{range .Model.ToOneRelationships -}}
	{{- $txt := txtsFromOneToOne $dot.Models $dot.Model . -}}
	{{$txt.Function.NameGo}} *{{$txt.ForeignModel.NameGo}}
	{{end -}}

	{{range .Model.ToManyRelationships -}}
	{{- $txt := txtsFromToMany $dot.Models $dot.Model . -}}
	{{$txt.Function.NameGo}} {{$txt.ForeignModel.Slice}}
	{{end -}}{{/* range tomany */}}
}

// {{$modelNameCamel}}L is where Load methods for each relationship are stored.
type {{$modelNameCamel}}L struct{}
{{end -}}
