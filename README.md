# SQLBunny

sqlbunny is a Go ORM based on code generation.

## Features
- Statically typed, fast generated code. No `interface{}`!
- Postgres fully supported, MySQL and others coming soon.
- Automatic migration generation (diffing the current migrations with the defined models)
- Relationship helper functions 
- Enums
- Structs. Can be reused across models, and are "flattened" to multiple SQL columns.
- Support for custom Go types in fields.

## Documentation

Documentation and guides are available at https://sqlbunny.io/

## Sneak peek

Models are defined like this:

```go
Model("book",
    Field("id", "string", PrimaryKey),
    Field("created_at", "time"),
    Field("name", "string"),
)
```

After running sqlbunny, a `Book` Go struct is generated, with functions to load and store objects from the database:

```go
// Insert a new book
book := &models.Book{
    ID:        "574389527",
    Name:      "Harry Potter and the Philosopher's Stone",
    CreatedAt: time.Now(),
}
err = book.Insert(ctx)

// Fetch all books
books, err := models.Books().All(ctx)
```

Want more? Check out the [Getting Started](https://sqlbunny.io/getting-started.html) guide.

## Credits

This project started out as a fork of the excellent [SQLBoiler](https://github.com/volatiletech/sqlboiler).
