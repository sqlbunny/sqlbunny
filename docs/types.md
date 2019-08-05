# Types 

Types are defined using `Type` in the top-level config (NOT inside a model declaration). Defining a type makes it available for use in `Field` definitions.

```go
Run(
    Type("my_type", typedefinition),
    Model("my_model", 
        Field("my_column", "my_type"),
    ),
)
```

`typedefinition` can be any of the type definitions explained below, or even from a plugin.

## Base types

The simplest type is `BaseType`, which maps to a Go type and an SQL type.

The Go type is specified using `Go` and `GoNull`, for non-nullable and nullable columns respectively. Types that require imports should be specified separating the import path and type name with a dot: `the/import/path.TypeName`.

Go structs should implement [`driver.Valuer`](https://golang.org/pkg/database/sql/driver/#Valuer) and [`sql.Scanner`](https://golang.org/pkg/database/sql/#Scanner), so the SQL driver knows how to marshal it to/from the datbase.

The database types must specify the type and the SQL value corresponding to the Go zero value to be used in [non-nullable column defaults](/models.html#default-values). Specifying the zero value is optional, but strongly encouraged.

```go
Type("string", BaseType{
    Go:       "string",
    GoNull:   "github.com/sqlbunny/sqlbunny/types/null.String",
    Postgres: SQLType{
        Type: "text",
        ZeroValue: "''",
    },
}),

Type("amount", BaseType{
    Go:       "github.com/exampleorg/exampleproject/types.Amount",
    GoNull:   "github.com/exampleorg/exampleproject/types.NullAmount",
    Postgres: SQLType{
        Type: "bigint",
        ZeroValue: "0",
    },
}),
```

Base types do not cause any code to be generated. 

## Enums

`Enum` is an integer-based enum. Every enum option is assigend an integer, starting with 0 for the first option. Therefore, the default/zero value of the enum type is the first option.

```go
Type("operation_state", Enum(
    "queued",      // 0
    "processing",  // 1
    "completed",   // 2
    "failed",      // 3
)),
```

:::warning
Since options are assigned values based on their order, you can't reorder, remove or add options (other than at the end). Doing so will break the data currently in the database, because the values are stored as their integers.
:::

:::tip
On the plus side, you can rename an existing option without breaking the data currently in the database.
:::

An enum generates a Go type in the models package:

```go
// The type is named after the enum.
var state models.OperationState

// Option constants are available under the enum name pluralized
state = models.OperationStates.Processing

// Convert to string
state.String()  // returns "processing"

// Convert from string
// If not a correct option, returns a bunny.InvalidEnumError
state, err = OperationStateFromString("completed")

// MarshalText, UnmarshalText are generated so the enum is
// marshaled as a string, option value in JSON & co.
```

A second type `NullOperationState` is also generated for use in nullable columns.

## Structs 

A struct generates a Go struct with the contained fields. When used as a type of a model field, they expand the field into multiple SQL columns, one per struct field.

Structs can be nested in other structs.

Primary keys, indexes, uniques, and foreign keys defined on struct fields are "propagated" to models using the struct, as if they were defined on the model directly.

You can define primary keys, indexes, uniques, and foreign keys in a model that involve inner struct fields specifying them with a dot `.`.

```go
Type("money", Struct(
    Field("amount", "int64"),
    Field("currency", "currency", Index),
)),
Model("account",
    Field("id", "string", PrimaryKey),
    Field("balance", "money"),
    Index("balance.amount"),
),
```

The above struct and model would generate the following Go code:

```go
type Money struct {
	Amount   int64    `bunny:"amount"`
	Currency Currency `bunny:"currency"`
}
type Account struct {
    ID       string  `bunny:"id"`
    Balance  Value   `bunny:"balance__,bind"`
}
```

and the following SQL tables:

```sql
CREATE TABLE "account" (
    "id" text NOT NULL DEFAULT '',
    "balance__amount" bigint NOT NULL DEFAULT 0,
    "balance__currency" integer NOT NULL DEFAULT 0
);

-- index defined from inside the struct
CREATE INDEX CONCURRENTLY "account___balance__currency___idx" ON "account" ("balance__currency");

-- index involving a struct field, defined outside the struct
CREATE INDEX CONCURRENTLY "account___balance__amount___idx" ON "account" ("balance__amount");

ALTER TABLE "account"
    ADD CONSTRAINT "account_pkey" PRIMARY KEY ("id");
```


As you can see, the field `balance` has been expanded into two SQL columns, `balance__amount` and `balance__currency`. The column names are generated joining the model field and struct field names with `__`. (The reason for using double underscore is a single underscore might clash with another model field.)

## Nullable structs

Fields of struct types support the `Null` option. This makes the struct as a whole optional (so either all of the fields are set, or none is).

For example, you're building a service that allows users to optionally verify their real-world identity with a government-issued identity document. You want a User object to have all document fields set, or none.

With nullable structs, this is easy:

```go
Type("identity_document", Struct(
    Field("number", "string"),
    Field("expiration_date", "date"),
    Field("country", "country"),
)),
Model("user",
    Field("id", "string", PrimaryKey),
    Field("identity_document", "identity_document", Null),
),
```

This generates the following Go code:

```go
type IdentityDocument struct {
	Number         string         `bunny:"number"`
	ExpirationDate _import00.Date `bunny:"expiration_date"`
	Country        Country        `bunny:"country"`
}
type NullIdentityDocument struct {
	IdentityDocument IdentityDocument
	Valid            bool
}
type User struct {
	ID               string               `bunny:"id"`
	IdentityDocument NullIdentityDocument `bunny:"identity_document__,bind,null:identity_document"`
}
```

and the following SQL tables:

```sql
CREATE TABLE "user" (
    "id" text NOT NULL DEFAULT '',
    "identity_document__number" text,
    "identity_document__expiration_date" date,
    "identity_document__country" integer,
    "identity_document" boolean NOT NULL DEFAULT false
);
```

The struct nullability is achieved by making all struct columns nullable (even if they were not nullable in the struct), and adding an extra bool column named after the struct. If this column is `true`, the `identity_document` field is set. If `false`, the field is not set (and the other columns should be NULL).

:::warning
At any moment, the struct field columns should be NULL **if and only if** the boolean column is false. 

This is currently NOT enforced through SQL constraints. This will be adressed in a future release.

The generated code preserves this invariant when inserting/updating full objects, but it can be violated by isnerting/updating specifying a partial column list. The best way to ensure sanity is to always include either all or no struct columns in the column list.
:::