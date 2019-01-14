# Querying

To run a query, you have to [start](#starters) it, specify some [query mods](#query-mods), and then execute the desired action using one of the available [finishers](#finishers).

## Starters

Starter methods are generated for each model, named after the plural. For example, `models.Books()`.

```go
models.Books(qm.Where("id=?", 5)).All(ctx)
```

You can also use `models.NewQuery()`. This allows you to specify the table name yourself, to do more advanced queries. Queries created this way can't use finishers that return model instances (such as `All` and `One`).

```go
models.NewQuery(
    qm.Where("id=?", 5)
).Bind(ctx, &objs)
```

## Query mods

### Select
`Select` can be used to specify which columns to select. 
```go
Select("id", "name")
```

### From
`From` specifies from which table to select. It can be used to give an alias to the table, to help with complex queries.
```go
From("book AS b")
```

### Where, WhereIn
`Where` specifies a condition used to filter rows in the `WHERE` clause. Multiple `Where` mods in the same query will be ANDed together.

```go
Where("name=?", "Daniel")
Where("name=? OR name=?", "Daniel", "John")
Where("age >= ? AND age <= ?", 18, 24)
```

`WhereIn` generates a `WHERE ... IN (...)` clause. It is helpful to deal with variable length conditions, because the `?` is expanded to the correct number of placeholders.

```go
WhereIn("weight in ?", 84)
WhereIn("height in ?", 183, 177, 204)
WhereIn("name, age in ?", "John", 24, "Tim", 33) // Generates: WHERE ("name","age") IN (($1,$2),($3,$4))
```

### Limit, Offset
`Limit` and `Offset` set the SQL `LIMIT` and `OFFSET` clause. 

```go
Limit(25), Offset(100),
```

### OrderBy

```go
OrderBy("age")
OrderBy("age, height")
OrderBy("age DESC, height ASC")
OrderBy("user_id ASC NULLS FIRST")
```

### GroupBy

```go
GroupBy("age")
GroupBy("age, gender")
```

### Having

```go
Having("COUNT(*) > ?", 10)
```

### For
`For` sets the SQL `FOR` clause, using for controlling locking behavior in transactions. 

```go
For("update nowait")
```

### InnerJoin
```go
InnerJoin("user u ON u.id = p.user_id")
InnerJoin("post p ON p.user_id = u.id AND p.type = ?", models.PostTypes.Story)
```

### Load

`Load` eagerly loads the models in the referred relationships. See [Eager loading](/relationships#eager-loading).

### SQL

`SQL` overrides the entire SQL query. When present, all other query mods will be ignored.
```go
models.Books(qm.SQL("SELECT * FROM book WHERE id=$1", 1234))
```

## Finishers

All finishers take as a first argument a `context.Context` instance. This context must contain a reference to the DB to use, using `bunny.ContextWithDB`.

### All

Returns all the matched rows in a slice. If no rows match, an empty slice is returned (it's not an error).

```go
books, err = models.Books(
    qm.Where("genre = ?", models.BookGenres.Fiction),
).All(ctx)
```

### One

Returns one matched row. If no rows match, `sql.ErrNoRows` is returned. If more than one row matches, `bunny.ErrMultipleRows` is returned.

```go
books, err = models.Books(
    qm.Where("title = ?", "Harry Potter and the Philosopher's Stone"),
).One(ctx)
```

### First

Returns the first matched row. If no rows match, `sql.ErrNoRows` is returned. If more than one row matches, the first row is returned (and no error is returned). The first row is defined by the query order, so if you do care about it, do specify `qm.OrderBy`.

It is recommended to use `One` instead, unless you specifically want to tolerate multiple rows.

```go
books, err = models.Books(
    qm.Where("title = ?", "Harry Potter and the Philosopher's Stone"),
).First(ctx)
```

### Count

Returns the integer count of rows matching the mods.

```go
count, err = models.Books(
    qm.Where("genre = ?", models.BookGenres.Fiction),
).Count(ctx)
```

### Exists

Returns a boolean specifying whether at least one row exists that matches the mods.

```go
exists, err = models.Books(
    qm.Where("title = ?", "Harry Potter and the Philosopher's Stone"),
).Exists(ctx)
```

### UpdateMapAll
Updates all rows matching the mods with the specified field values.

```go
err := models.Users(
    qm.Where("age >= ?", 18),
).UpdateMapAll(ctx, models.M{"is_legal_age": true})
```

### DeleteAll
Deletes all rows matching the mods.
```go
err := models.Users(
    qm.Where("age >= ?", 18),
).DeleteAll(ctx)
```

### Bind
Executes the query, binding the result to the given struct or slice. See [binding](/query-bind.html) for more info.

### Query, QueryRow, Exec
Executes the query and returns the raw SQL result. These functions simply mirror thir [`sql.DB`](https://golang.org/pkg/database/sql/#DB) counterparts.

## Model functions

### Find

Convenience function to fetch a model instance by its primary key. Returns `sql.ErrNoRows` if not found.

```go
user, err := models.FindUser(ctx, 12) // Finds the user with ID 12 

// Equivalent to:
user, err := models.Users(qm.Where("id=?", 12)).One()
```

### Exists

Convenience function to check if a model with the given primary key exists.

```go
exists, err := models.UserExists(ctx, 12) // Does the user with ID 12 exist?

// Equivalent to:
exists, err := models.Users(qm.Where("id=?", 12)).Exists()
```

### Insert

`Insert` inserts a model instance to the database.

```go
book := &models.Book{
    ID:        "574389527",
    Name:      "Harry Potter and the Philosopher's Stone",
    CreatedAt: time.Now(),
}
err = book.Insert(ctx)
```

You can explicitly specify a list of columns to be inserted. If not set, all the model columns will be inserted.

```go
book := &models.Book{
    ID:        "574389527",
    CreatedAt: time.Now(),
}
err = book.Insert(ctx,
    models.BookColumns.ID,
    models.BookColumns.CreatedAt,
)
```

### Update

`Update` updates a model instance in the database.

```go
book.Name = "SuperNiceBook"
err = book.Update(ctx)
```

You can explicitly specify a list of columns to be updated. If not set, all non-primary-key columns will be updated.

```go
book.Name = "SuperNiceBook"
err = book.Insert(ctx,
    models.BookColumns.Name,
)
```

### Delete

Deletes the model from the database.

```go
err = book.Delete(ctx)
```

### Reload

If changes happen in the database and the data in the Go struct becomes outdated, you can reload it from the database using `Reload`.

```go
err = book.Reload(ctx)
```

## Slice functions

### UpdateMapAll
Updates all rows matching the mods with the specified field values.

```go
err = users.UpdateMapAll(ctx, models.M{"is_legal_age": true})
```

### ReloadAll
Reloads all the model instances in the slice from the database.

```go
err = users.ReloadAll(ctx)
```

### DeleteAll
Deletes all the objects in the slice.

```go
err = users.DeleteAll(ctx)
```
