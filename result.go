package gofp

import (
	"errors"
	"fmt"

	"github.com/Alsond5/gofp/either"
	"github.com/Alsond5/gofp/tuple"
)

type Result[T any] struct {
	value T
	err   error
	ok    bool
}

func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, ok: true}
}

func Err[T any](err error) Result[T] {
	if err == nil {
		var zero T
		return Result[T]{value: zero, ok: true}
	}

	return Result[T]{err: err, ok: false}
}

func Do(err error) Result[Unit] {
	if err != nil {
		return Err[Unit](err)
	}

	return Ok(Unit{})
}

func Of[T any](value T, err error) Result[T] {
	if err != nil {
		return Result[T]{err: err, ok: false}
	}

	return Result[T]{value: value, ok: true}
}

func Of2[A, B any](a A, b B, err error) Result[tuple.Pair[A, B]] {
	if err != nil {
		return Err[tuple.Pair[A, B]](err)
	}

	return Ok(tuple.Pair[A, B]{First: a, Second: b})
}

func Of3[A, B, C any](a A, b B, c C, err error) Result[tuple.Triple[A, B, C]] {
	if err != nil {
		return Err[tuple.Triple[A, B, C]](err)
	}

	return Ok(tuple.Triple[A, B, C]{First: a, Second: b, Third: c})
}

func Try[T any](f func() T) (r Result[T]) {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}

		if rp, ok := rec.(resultPanic); ok {
			r = Err[T](rp.err)
			return
		}

		panic(rec)
	}()

	return Ok(f())
}

func TryCatch[L, R any](try func() L, catch func(error) R) (e either.Either[L, R]) {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}

		if rp, ok := rec.(resultPanic); ok {
			e = either.Right[L](catch(rp.err))
			return
		}

		panic(rec)
	}()

	return either.Left[L, R](try())
}

func (r Result[T]) Unwrap() T {
	if !r.ok {
		panic(resultPanic{err: r.err})
	}

	return r.value
}

func (r Result[T]) UnwrapErr() error {
	if r.ok {
		panic(r.value)
	}

	return r.err
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
	if !r.ok {
		return defaultValue
	}

	return r.value
}

func (r Result[T]) UnwrapOrElse(f func(error) T) T {
	if !r.ok {
		return f(r.err)
	}

	return r.value
}

func (r Result[T]) UnwrapOrZero() T {
	if !r.ok {
		var zero T
		return zero
	}

	return r.value
}

func (r Result[T]) Expect(msg string) T {
	if !r.ok {
		panic(msg + ": " + r.err.Error())
	}

	return r.value
}

func (r Result[T]) ExpectErr(msg string) error {
	if r.ok {
		panic(fmt.Sprintf("%s: %v", msg, r.value))
	}

	return r.err
}

func (r Result[T]) IntoErr() error {
	return r.err
}

func (r Result[T]) IntoOk() T {
	return r.value
}

func (r Result[T]) Ok() Option[T] {
	if r.ok {
		return Some(r.value)
	}

	return None[T]()
}

func (r Result[T]) Err() Option[error] {
	if !r.ok {
		return Some(r.err)
	}

	return None[error]()
}

func (r Result[T]) Unpack() (T, error) {
	return r.value, r.err
}

func (r Result[T]) IsOk() bool {
	return r.ok
}

func (r Result[T]) IsOkAnd(f func(T) bool) bool {
	return r.ok && f(r.value)
}

func (r Result[T]) IsErr() bool {
	return !r.ok
}

func (r Result[T]) IsErrAnd(f func(error) bool) bool {
	return !r.ok && f(r.err)
}

func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if !r.ok {
		return Err[T](f(r.err))
	}

	return r
}

func (r Result[T]) OrElse(f func(error) Result[T]) Result[T] {
	if !r.ok {
		return f(r.err)
	}

	return r
}

func (r Result[T]) Or(alternative Result[T]) Result[T] {
	if !r.ok {
		return alternative
	}

	return r
}

func (r Result[T]) ContainsErr(target error) bool {
	if r.ok {
		return false
	}

	return errors.Is(r.err, target)
}

func (r Result[T]) IfOk(f func(T)) Result[T] {
	if r.ok {
		f(r.value)
	}

	return r
}

func (r Result[T]) IfErr(f func(error)) Result[T] {
	if !r.ok {
		f(r.err)
	}

	return r
}

func (r Result[T]) Tap(okFn func(T), errFn func(error)) Result[T] {
	if r.ok {
		if okFn != nil {
			okFn(r.value)
		}
	} else {
		if errFn != nil {
			errFn(r.err)
		}
	}

	return r
}
