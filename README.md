<h1 align="center">
  <a href="https://github.com/Alsond5/gofp">
    <picture>
      <img height="125" alt="gofp" src="https://raw.githubusercontent.com/Alsond5/gofp/main/.github/assets/logo.svg">
    </picture>
  </a>
  <br>
  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg">
  </a>
</h1>

<p align="center">
  <em><b>Gofp</b> is a <a href="https://doc.rust-lang.org/std/result/">Rust</a> inspired <b>functional primitives</b> library for <a href="https://go.dev/doc/">Go</a>. It's bringing type-safe `Result`, `Option`, and `Either` types with <b>zero allocations</b> and composable <b>error handling</b>.</em>
</p>

---

## Why

Go's `if err != nil` pattern is honest and explicit. But it doesn't compose. You can't chain transformations, propagate errors through a pipeline or express "this value might not exist" in the type system.

`gofp` fills that gap not by hiding errors but by making them first-class values you can map, chain, and transform.

## Install

```bash
go get github.com/Alsond5/gofp
```

Requires **Go 1.23+**.

## Packages

| Package | Description |
|---|---|
| `gofp` | `Result[T]`, `Option[T]`, `Unit` |
| `gofp/result` | `Map`, `FlatMap`, `AllOf`, `Partition`, `FirstOk` Result transformations and combinators |
| `gofp/option` | `Map`, `FlatMap`, `Zip`, `Match` Option transformations and combinators |
| `gofp/either` | `Either[L,R]` two-outcome type for domain branching |
| `gofp/must` | Panic helpers for initialization |

## Result\[T\]

Represents either a successful value (`Ok`) or an error (`Err`). Drop-in replacement for Go's `(T, error)` pattern.

```go
// Construct
r := gofp.Ok(42)
r := gofp.Err[int](errors.New("something went wrong"))
r := gofp.Of(strconv.Atoi("42"))        // bridges (T, error) → Result[T]

// Check
r.IsOk()
r.IsErr()

// Unwrap
r.Unwrap()                              // panics on Err
r.UnwrapOr(0)                           // fallback value
r.UnwrapOrElse(func(err error) int { return -1 })

// Transform — free functions, type-safe T→U
result.Map(r, func(n int) string { return strconv.Itoa(n) })    // Result[string]
result.FlatMap(r, func(n int) gofp.Result[string] { ... })      // Result[string]

// Chain multiple Results
result.FlatMap(
    gofp.ResultFlatMap(parse(s), validate),
    process,
)

// Side effects
r.IfOk(func(n int) { log.Println(n) }).
    IfErr(func(err error) { log.Println(err) })

// Combine
result.AllOf(r1, r2, r3)                 // Result[[]T] — first Err wins
result.All2(r1, r2)                      // Result[Pair[A, B]]
result.Partition(results...)             // ([]T, []error)
result.FirstOk(r1, r2, r3)               // first Ok, or all errors joined
```

### Try — Go's answer to `?`

```go
func handle(s string) gofp.Result[string] {
    return gofp.Try(func() string {
        age := parse(s).Unwrap()        // Err → caught by Try
        age  = validate(age).Unwrap()   // Err → caught by Try
        return format(age).Unwrap()     // Err → caught by Try
    })
    // real panics (nil pointer, index out of range) are re-panicked
}
```

## Option\[T\]

Represents a value that may or may not exist. Replaces `nil` checks and pointer abuse.

```go
// Construct
o := gofp.Some(42)
o := gofp.None[int]()
o := gofp.FromPtr(ptr)                  // nil ptr → None
o := gofp.FromZero(val)                 // zero value → None

// Check
o.IsSome()
o.IsNone()

// Unwrap
o.Unwrap()                              // panics on None
o.UnwrapOr(0)
o.UnwrapOrElse(func() int { return compute() })
o.UnwrapOrZero()                        // returns zero value

// Transform
o.Filter(func(n int) bool { return n > 0 })

option.Map(o, func(n int) string { ... })   // Option[string]
option.FlatMap(func(n int) gofp.Option[string] { ... })

// Boolean combinators
o.Or(other)        // first Some
o.Xor(other)       // Some only if exactly one is Some
option.And(other)  // Some only if both Some

// Pair operations
option.Zip(a, b)                    // Option[Pair[A, B]]
option.ZipWith(a, b, func(x, y) z)  // Option[Z]

// Side effects
o.IfSome(func(n int) { ... }).
 IfNone(func() { ... })

o.Match(
    func(n int) { ... },                // Some branch
    func() { ... },                     // None branch
)

// Convert to Result
o.OkOr(ErrNotFound)
o.OkOrElse(func() error { return ErrNotFound })
```

## Either\[L, R\]

A value that is either `Left(L)` or `Right(R)`. Unlike `Result`, neither side implies failure both are valid domain values.

```go
// Construct
e := either.Left[int, string](42)
e := either.Right[int, string]("guest")

// Check
e.IsLeft()
e.IsRight()

// Fold — primary way to consume an Either
msg := either.Fold(e,
    func(n int) string { return fmt.Sprintf("number: %d", n) },
    func(s string) string { return "string: " + s },
)

// Transform
either.MapLeft(e, func(n int) int { return n * 2 })
either.MapRight(e, func(s string) string { return strings.ToUpper(s) })
either.MapBoth(e, leftFn, rightFn)

// Side effects
e.IfLeft(func(n int) { log.Println(n) }).
 IfRight(func(s string) { log.Println(s) })

// Swap sides
e.Swap()                                // Either[R, L]

// Partition a slice
lefts, rights := either.Partition(eithers)

// Merge when both sides are the same type
either.Merge(either.Left[int, int](42)) // → 42
```

### TryCatch — Try with two return types

```go
e := either.TryCatch(
    func() int {                        // try  → Left on success
        age := parse(s).Unwrap()
        return validate(age).Unwrap()
    },
    func(err error) string {            // catch → Right on error
        return "invalid: " + err.Error()
    },
)
// Either[int, string]

either.Fold(e,
    func(age int) { fmt.Println("age:", age) },
    func(msg string) { fmt.Println("error:", msg) },
)
```

Real panics (nil pointer, index out of range) are **re-panicked**, not swallowed.

## must

Panic helpers for program initialization. **Not for request handling.**

```go
import "github.com/Alsond5/gofp/must"

// stdlib bridge
db := must.Do(sql.Open("postgres", dsn))

// invariants
must.Be(port > 0 && port < 65536, "invalid port")
must.Bef(len(items) > 0, "expected at least %d items, got %d", 1, len(items))

// environment
dsn  := must.Env("DATABASE_URL")           // panics if missing
port := must.EnvIntOr("PORT", 8080)        // fallback, panics if set but invalid

// parsing
re   := must.Regexp(`^[^@]+@[^@]+\.[^@]+$`)
tmpl := must.Template(template.ParseFiles("email.html"))

// nil safety
cfg := must.NotNil(loadConfig(), "config must not be nil")

// collection access
first := must.Index(items, 0, "items is empty")
val   := must.Key(m, "key", "missing required key")

// unreachable
default: must.Never("unhandled case")
```

## Example

```go
var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidEmail = errors.New("invalid email")
)

func getUser(db map[uint64]User, id uint64) gofp.Result[User] {
    return gofp.FromZero(db[id]).OkOr(ErrNotFound)
}

func validateEmail(email string) gofp.Result[string] {
    if strings.Contains(email, "@") {
        return gofp.Ok(email)
    }
    return gofp.Err[string](ErrInvalidEmail)
}

func domainAccess(email string) gofp.Option[uint8] {
    switch {
    case strings.HasSuffix(email, "@admin.com"):
        return gofp.Some[uint8](10)
    case strings.HasSuffix(email, "@user.com"):
        return gofp.Some[uint8](1)
    default:
        return gofp.None[uint8]()
    }
}

func computeAccess(db map[uint64]User, id uint64) gofp.Result[gofp.Option[uint8]] {
    return result.FlatMap(
        getUser(db, id),
        func(user User) gofp.Result[gofp.Option[uint8]] {
            email := user.Email
            return option.MapOr(
                email,
                gofp.Ok(gofp.None[uint8]()),
                func(email string) gofp.Result[gofp.Option[uint8]] {
                    return result.Map(validateEmail(email), domainAccess)
                },
            )
        },
    )
}

func main() {
    computeAccess(db, 1).
        IfOk(func(opt gofp.Option[uint8]) {
            opt.Match(
                func(level uint8) { fmt.Println("access level:", level) },
                func()            { fmt.Println("no domain match") },
            )
        }).
        IfErr(func(err error) { fmt.Println("error:", err) })
}
```

---

## Design Decisions

**Why `ok bool` in structs?**
Go structs always have a zero value. Without an explicit flag, `Result[int]{}` and `Ok(0)` are structurally identical indistinguishable. The `ok`/`some` field makes the zero value explicitly invalid.

**Why value receivers is not pointer receivers?**
Value receivers keep `Result[T]` and `Option[T]` on the stack. Pointer receivers cause heap allocation via escape analysis. For types used in hot paths, this matters.

**Why free functions for `Map`, `FlatMap` etc. ?**
Go methods cannot introduce new type parameters. `func (r Result[T]) Map[U](f func(T) U) Result[U]` is illegal. Free functions are the only type-safe way to express `T → U` transformations.

**Why not `pipe`?**
A type-safe heterogeneous pipeline (`T → U → V`) requires methods that introduce new type parameters which Go doesn't allow. `T → T` pipes exist but add little over plain function calls. `FlatMap` chains cover the real use cases.

## License
gofp is licensed under the [MIT License](./LICENSE).
