# Query binding

The binding system lies at the heart of all sqlbunny queries. It is in charge of converting SQL rows into Go structs.

It matches SQL columns to Go struct fields, and copies their values over.

Binding behavior is controlled by the `bunny` struct tag.

## Field binds

To bind a field to an SQL column, specify its SQL column name in the `bunny` tag. Fields without a `bunny` tag are never bound.

SQL-to-Go conversion is done using the same rules as [`sql.Rows.Scan()`](https://golang.org/pkg/database/sql/#Rows.Scan). In a nutshell, it means destination types should be either supported directly by the SQL driver, or implement [`sql.Scanner`](https://golang.org/pkg/database/sql/#Scanner).

```go
type MyStruct struct {
    Foo     int    `bunny:"foo"`
    Bar     string `bunny:"bar"`
    Ignored string
}
var v MyStruct
err = queries.Raw("SELECT 1 AS foo, 'fun' AS bar, 'hello' AS ignored").Bind(ctx, &v)
// v contains foo=1, bar="fun", and ignored=""
```

## Recursive binds

If the ",bind" option is specified on a field of struct type, binding will recurse into it
to look for fields for binding. "name" is added as a prefix to the SQL column names
of the inner fields.

```go
type Bar struct {
    One int `bunny:"one"`
    Two int `bunny:"two"`
}
type MyStruct struct {
    Foo     int `bunny:"foo"`
    Bar     Bar `bunny:"bar__,bind"`
}
var v MyStruct
err = queries.Raw("SELECT 1 AS foo, 2 AS bar__one, 3 as bar__two").Bind(ctx, &v)
// v contains foo=1, bar={ one=2, two=3 }
```

## Aggregating

Bind to a custom struct to read the results of an aggregation.
```go
type BooksByYear struct {
    Year  int `bunny:"year"`
    Count int `bunny:"count"`
}

var years []BooksByYear

err = models.Books(
    qm.Select("year", "count(*) AS count"),
    qm.GroupBy("year"),
    qm.OrderBy("year ASC"),
).Bind(ctx, &years)
```

## Joining

Binding to a custom "join struct" can be handy for receiving the results of a join.

```go
type BookAndAuthor struct {
    models.Author `bunny:"author.,bind"`
    models.Book   `bunny:"book.,bind"`
}

var res []BookAndAuthor

err = models.Authors(
    qm.Select(
        "author.id", "author.name",
        "book.id", "book.author_id", "book.title", "book.year",
    ),
    qm.InnerJoin("book ON book.author_id = author.id"),
).Bind(ctx, &res)

// This wil execute the following query:
// SELECT 
//     "author"."id" as "author.id", "author"."name" as "author.name",
//     "book"."id" as "book.id", "book"."author_id" as "book.author_id", "book"."title" as "book.title", "book"."year" as "book.year"
// FROM "author"
// INNER JOIN book ON book.author_id = author.id;
```

## Extending models

Recursive binding with empty prefix can be used to "extend" a model query with an extra computed column.

You simply define your own struct extending `Author`, and bind to it. Since `Author` is recursively bound with empty prefix, all the `Author` 

```go
type AuthorWithCount struct {
    models.Author `bunny:",bind"`
    BookCount     int `bunny:"book_count"`
}

var authors []AuthorWithCount

err = models.Authors(
    qm.Select("*", "(SELECT COUNT(*) FROM book WHERE book.author_id = author.id) AS book_count"),
).Bind(ctx, &authors)
```

This gives us all the authors with their respective book counts!
