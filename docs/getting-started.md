# Getting Started

## Installation

```
go get -u github.com/kernelpayments/sqlbunny
```
## Project structure

sqlbunny's code generation is a library, not a command-line tool.

The recommended usage is making your own executable binary (ie, a Go `main` package) that calls the sqlbunny library. 

:::tip
This binary should be completely separate from your main application binary.

Having both functionalities in a single binary might seem simpler, but it's not a good 
idea! You wouldn't be able to use sqlbunny if the application
code does not compile (due to, for example, not having generated the sqlbunny models yet!).
:::

An example folder structure looks like this. To run code generation, run the `sqlbunny` executable. To run your app, run the `app` executable.

```
- cmd/
    - sqlbunny/      <- Runs sqlbunny, with your own configuration
        - main.go
    - app/           <- Runs your application
        - main.go
- models/            <- autogenerated models package
    - book.go        
    - bunny_*.go
```

## Defining models

Put the following code at `./cmd/sqlbunny/main.go`:

```go
package main

import (
	. "github.com/kernelpayments/sqlbunny/gen/core"
	"github.com/kernelpayments/sqlbunny/gen/migration"
	"github.com/kernelpayments/sqlbunny/gen/typelib"
)

func main() {
	Run(
		&typelib.Plugin{},
		&migration.Plugin{},

		Model("book",
			Field("id", "string", PrimaryKey),
			Field("created_at", "time"),
			Field("name", "string"),
		),
	)
}
```

The main function calls sqlbunny's `Run`, which takes a slice of configuration items. Everything is done through config items: enabling plugins, defining types and models...

In this case, the first two items enable the `typelib` and `migration` plugins. `typelib` defines many useful data types, and `migration` adds support for generating migrations.

The last configuration item defines a model named `book`.

## Generating

To generate the models, run the sqlbunny package:

```
go run ./cmd/sqlbunny/main.go gen
```

This will generate the `models` package. Let's take a look at it:
```
$ ls models/
bunny_array_utils.go  book.go  bunny_queries.go  bunny_table_names.go  bunny_types.go
```

sqlbunny has created `book.go` with our model, along with a few utility files named `bunny_*.go`. You can check out the `Book` struct definition:

```go
type Book struct {
	ID        string    `bunny:"id" json:"id" `
	CreatedAt time.Time `bunny:"created_at" json:"created_at" `
	Name      string    `bunny:"name" json:"name" `
	R         *bookR    `json:"-" toml:"-" yaml:"-"`
	L         bookL     `json:"-" toml:"-" yaml:"-"`
}
```

As you can see, the struct has the 3 fields we defined previously, and two extra `R` and `L` fields, used for [eager loading](TODO).

`book.go` also contains many functions to query and modify `Book` objects in the database.

## Creating the tables

We have the Go code to interact with a database, but we don't have such a database with the correct tables yet. Let's fix that.

Run the following:
```
go run ./cmd/sqlbunny/main.go gensql
```

This converts your model definitions to a series of SQL statements. Run them in your database to create the `book` table.

```sql
CREATE TABLE "book" (
    "id" text NOT NULL DEFAULT '',
    "created_at" timestamptz NOT NULL,
    "name" text NOT NULL DEFAULT ''
);

ALTER TABLE "book"
    ADD CONSTRAINT "book_pkey" PRIMARY KEY ("id");
```

:::tip
The SQL generated by `gensql` is only suitable for creating all the tables in an empty database. If you make changes to your model
definitions and want to incrementally migrate existing databases, you will have to properly configure the [migrations](TODO) plugin.

Additionally, sqlbunny doesn't force you to use the `migrations` plugin. You can use any tool to manage your schema. As long as the table and column names match, everything will work!
:::

## Using the generated models

Now that you have a `models` package, let's do some queries! All the code below goes in `./cmd/app/main.go`.

First, we'll need some imports, and the main function

```go
package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/kernelpayments/sqlbunny/runtime/bunny"
	"example.com/sqlbunny_demo/models"
)

func main() {
    ...
}
```

All ORM functions need a context with an active database. To set the database in a context, wrap the context with `bunny.ContextWithDB`. The context can also be used to cancel in-progress queries (it's passed through to the SQL driver).

```go
db, err := sql.Open("postgres", "host=localhost port=5432 dbname=postgres user=postgres password=postgres sslmode=disable")
if err != nil {
	panic(err)
}

ctx := context.Background()
ctx = bunny.ContextWithDB(ctx, db)
```

Once you have a context, you can start doing ORM operations, such as inserting a book:

```go
// Insert a new book
book := &models.Book{
    ID:        "574389527",
    Name:      "Harry Potter and the Philosopher's Stone",
    CreatedAt: time.Now(),
}
err = book.Insert(ctx)
if err != nil {
    // handle error
}
```

...or fetching all books.   

```go
// Fetch all books
books, err := models.Books().All(ctx)
if err != nil {
	panic(err)
}
```

You now have a working sqlbunny project. Enjoy!