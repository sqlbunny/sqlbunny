# Models

Models are defined using `Model` in the configuration.

```go
Model("book",
    Field("id", "string", PrimaryKey),
    Field("created_at", "time"),
    Field("name", "string"),
)
```

## Fields

Fields are defined using `Field`, specifying the name, the type, and a list of potential options.
```go
    Field("id", "string", options...),
```


## Nullable fields

Fields can be defined nullable using `Null`.

```go
Model("user",
    Field("id", "string", PrimaryKey),
    Field("created_at", "time"),
    Field("deleted_at", "time", Null),
)
```

## Default values

Setting field default values is intentionally NOT supported in sqlbunny. 

Columns get SQL defaults set to the Go zero value if not nullable, and no default if nullable (ie the default is `null`). This means `0` for numeric types, `""` for strings, etc. This matches the Go behavior when initializing struct fields.

:::tip
Consider what would happen in this case if default values were hypothetically supported:

```go
Model("user",
    Field("id", "string", PrimaryKey),
    Field("is_active", "bool", Default(true)),
)

u := &models.User{ID: "1234"}
u.Insert(ctx)
```

`Insert` would receive `is_active = false`, but it can't know if it's `false` because it was never set (so the default of `true` should be used) or it has been explicitly set to `false` (so the inserted row should have `false`).

Both choices can be surprising to users and lead to potential bugs. To avoid this, sqlbunny takes an opinionated stance, and makes SQL's behavior more closely match Go's.
:::

## Primary key

Models must have exactly one primary key. 

If the primary key contains one column, you can specify it in the `Field` declaration.

```go
    Field("id", "string", PrimaryKey),
```

If the primary key contains multiple columns, you have to specify it with its own declaration.

```go
    Field("user_name", "string"),
    Field("repo_name", "string"),
    PrimaryKey("user_name", "repo_name"),
```

## Indexes

Similarly to primary keys, indexes can be defined in the field declaration or with their own declaration.

```go
    Field("created_at", "time", Index),
```

```go
    Field("created_at", "time"),
    Field("user_id", "string"),
    Index("user_id", "created_at"),
```


## Unique constraints

Same as indexes.

```go
    Field("user_name", "string", Unique),
```

```go
    Field("user_name", "time"),
    Field("repo_name", "string"),
    Unique("user_name", "repo_name"),
```

## Foreign keys

You can declare a field to be a foreign key to another model. In this case, the field type must be equal to the destination model primary key type.

```go
    Field("user_id", "string", ForeignKey("user")),
```

Foreign keys to models with a multi-column primary key are also supported.

```go
Model("repo",
    Field("organization_id", "string"),
    Field("repo_id", "string"),
    Field("path", "string"),
    PrimaryKey("organization_id", "repo_id"),
    ...
),

Model("repo_file",
    Field("organization_id", "string"),
    Field("repo_id", "string"),
    ModelForeignKey("repo", "organization_id", "repo_id"),
    Field("path", "string"),
    PrimaryKey("organization_id", "repo_id", "path"),
),
```

Foreign keys spanning multiple columns (needed to refer to models with multi-column primary keys) are currently not supported.

## Field tags

Struct field tags can be specified with `Tag`.
```go
Field("user_name", "string", Tag("json", "userName"), Tag("foo", "bar")),
```
would generate the following struct field:
```go
    UserName string `bunny:"user_name" json:"userName" foo:"bar"`
```
