package option

import (
	"github.com/Alsond5/gofp"
	"github.com/Alsond5/gofp/tuple"
)

func And[T, U any](o gofp.Option[T], other gofp.Option[U]) gofp.Option[U] {
	if o.IsNone() {
		return gofp.None[U]()
	}

	return other
}

func AndThen[T, U any](o gofp.Option[T], f func(T) gofp.Option[U]) gofp.Option[U] {
	if o.IsNone() {
		return gofp.None[U]()
	}

	return f(o.Unwrap())
}

func Zip[T, U any](o gofp.Option[T], other gofp.Option[U]) gofp.Option[tuple.Pair[T, U]] {
	if o.IsSome() && other.IsSome() {
		return gofp.Some(tuple.Pair[T, U]{First: o.Unwrap(), Second: other.Unwrap()})
	}

	return gofp.None[tuple.Pair[T, U]]()
}

func ZipWith[T, U, R any](o gofp.Option[T], other gofp.Option[U], f func(T, U) R) gofp.Option[R] {
	if o.IsSome() && other.IsSome() {
		return gofp.Some(f(o.Unwrap(), other.Unwrap()))
	}

	return gofp.None[R]()
}

func Unzip[T, U any](o gofp.Option[tuple.Pair[T, U]]) (gofp.Option[T], gofp.Option[U]) {
	if o.IsSome() {
		p := o.Unwrap()
		return gofp.Some(p.First), gofp.Some(p.Second)
	}

	return gofp.None[T](), gofp.None[U]()
}

func Flatten[T any](o gofp.Option[gofp.Option[T]]) gofp.Option[T] {
	if o.IsNone() {
		return gofp.None[T]()
	}

	return o.Unwrap()
}

func Map[T, U any](o gofp.Option[T], f func(T) U) gofp.Option[U] {
	if o.IsNone() {
		return gofp.None[U]()
	}

	return gofp.Some(f(o.Unwrap()))
}

func MapOr[T, U any](o gofp.Option[T], defaultValue U, f func(T) U) U {
	if o.IsNone() {
		return defaultValue
	}

	return f(o.Unwrap())
}

func MapOrElse[T, U any](o gofp.Option[T], defaultFn func() U, f func(T) U) U {
	if o.IsNone() {
		return defaultFn()
	}

	return f(o.Unwrap())
}

func FlatMap[T, U any](o gofp.Option[T], f func(T) gofp.Option[U]) gofp.Option[U] {
	if o.IsNone() {
		return gofp.None[U]()
	}

	return f(o.Unwrap())
}

func Match[T, U any](o gofp.Option[T], someFn func(T) U, noneFn func() U) U {
	if o.IsSome() {
		return someFn(o.Unwrap())
	}

	return noneFn()
}
