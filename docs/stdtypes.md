# Standard types

The standard type library provides a set of common pre-defined types. Its usage is recommended for convenience, and for consistency across projects.

To use stdtypes, make sure the plugin is enabled in the sqlbunny configuration:

```go
Run(
    &stdtypes.Plugin{},
    ...
)
```

## Types 

sqlbunny | Go          | Postgres
---------|-------------|---------------
int16    | int16       | smallint
int32    | int32       | integer
int64    | int64       | bigint
float32  | float32     | real
float64  | float64     | double precision
bool     | bool        | boolean
string   | string      | text
bytea    | []byte      | bytea
jsonb    | types.JSON  | jsonb
time     | time.Time   | timestamptz
