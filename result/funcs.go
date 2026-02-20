package result

import (
	"errors"

	"github.com/Alsond5/gofp"
	"github.com/Alsond5/gofp/tuple"
)

func Map[T, U any](r gofp.Result[T], f func(T) U) gofp.Result[U] {
	if r.IsErr() {
		return gofp.Err[U](r.UnwrapErr())
	}

	return gofp.Ok(f(r.Unwrap()))
}

func FlatMap[T, U any](r gofp.Result[T], f func(T) gofp.Result[U]) gofp.Result[U] {
	if r.IsErr() {
		return gofp.Err[U](r.UnwrapErr())
	}

	return f(r.Unwrap())
}

func And[T, U any](r gofp.Result[T], other gofp.Result[U]) gofp.Result[U] {
	if r.IsErr() {
		return gofp.Err[U](r.UnwrapErr())
	}

	return other
}

func AndThen[T, U any](r gofp.Result[T], f func(T) gofp.Result[U]) gofp.Result[U] {
	if r.IsErr() {
		return gofp.Err[U](r.UnwrapErr())
	}

	return f(r.Unwrap())
}

func Flatten[T any](r gofp.Result[gofp.Result[T]]) gofp.Result[T] {
	if r.IsErr() {
		return gofp.Err[T](r.UnwrapErr())
	}

	return r.Unwrap()
}

func All2[A, B any](a gofp.Result[A], b gofp.Result[B]) gofp.Result[tuple.Pair[A, B]] {
	if a.IsErr() {
		return gofp.Err[tuple.Pair[A, B]](a.UnwrapErr())
	}
	if b.IsErr() {
		return gofp.Err[tuple.Pair[A, B]](b.UnwrapErr())
	}

	return gofp.Ok(tuple.Pair[A, B]{First: a.Unwrap(), Second: b.Unwrap()})
}

func All3[A, B, C any](a gofp.Result[A], b gofp.Result[B], c gofp.Result[C]) gofp.Result[tuple.Triple[A, B, C]] {
	if a.IsErr() {
		return gofp.Err[tuple.Triple[A, B, C]](a.UnwrapErr())
	}
	if b.IsErr() {
		return gofp.Err[tuple.Triple[A, B, C]](b.UnwrapErr())
	}
	if c.IsErr() {
		return gofp.Err[tuple.Triple[A, B, C]](c.UnwrapErr())
	}

	return gofp.Ok(tuple.Triple[A, B, C]{First: a.Unwrap(), Second: b.Unwrap(), Third: c.Unwrap()})
}

func AllOf[T any](results ...gofp.Result[T]) gofp.Result[[]T] {
	values := make([]T, 0, len(results))
	for _, r := range results {
		if r.IsErr() {
			return gofp.Err[[]T](r.UnwrapErr())
		}

		values = append(values, r.Unwrap())
	}

	return gofp.Ok(values)
}

func AllOfCollectErrs[T any](results ...gofp.Result[T]) gofp.Result[[]T] {
	values := make([]T, 0, len(results))

	var errs []error
	for _, r := range results {
		if r.IsErr() {
			errs = append(errs, r.UnwrapErr())
		} else {
			values = append(values, r.Unwrap())
		}
	}
	if len(errs) > 0 {
		return gofp.Err[[]T](errors.Join(errs...))
	}

	return gofp.Ok(values)
}

func FirstOk[T any](results ...gofp.Result[T]) gofp.Result[T] {
	if len(results) == 0 {
		return gofp.Err[T](gofp.ErrNoResults)
	}

	var errs []error
	for _, r := range results {
		if r.IsOk() {
			return r
		}

		errs = append(errs, r.UnwrapErr())
	}

	return gofp.Err[T](errors.Join(errs...))
}

func AllErrors[T any](results ...gofp.Result[T]) []error {
	var errs []error
	for _, r := range results {
		if r.IsErr() {
			errs = append(errs, r.UnwrapErr())
		}
	}

	return errs
}

func Partition[T any](results ...gofp.Result[T]) (values []T, errs []error) {
	for _, r := range results {
		if r.IsOk() {
			values = append(values, r.Unwrap())
		} else {
			errs = append(errs, r.UnwrapErr())
		}
	}

	return values, errs
}

func PartitionResults[T any](results ...gofp.Result[T]) (oks []gofp.Result[T], errs []gofp.Result[T]) {
	for _, r := range results {
		if r.IsOk() {
			oks = append(oks, r)
		} else {
			errs = append(errs, r)
		}
	}

	return oks, errs
}

func MapAll[T, U any](slice []T, f func(T) gofp.Result[U]) []gofp.Result[U] {
	out := make([]gofp.Result[U], len(slice))
	for i, v := range slice {
		out[i] = f(v)
	}

	return out
}

func MapAllOk[T, U any](slice []T, f func(T) gofp.Result[U]) gofp.Result[[]U] {
	return AllOf(MapAll(slice, f)...)
}
