package gofp

type Option[T any] struct {
	value T
	ok    bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{value: value, ok: true}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return None[T]()
	}

	return Some(*ptr)
}

func FromZero[T comparable](value T) Option[T] {
	var zero T
	if value == zero {
		return None[T]()
	}

	return Some(value)
}

func (o Option[T]) IsSome() bool { return o.ok }

func (o Option[T]) IsSomeAnd(f func(T) bool) bool {
	return o.ok && f(o.value)
}

func (o Option[T]) IsNone() bool { return !o.ok }

func (o Option[T]) IsNoneOr(f func(T) bool) bool {
	if !o.ok {
		return true
	}

	return f(o.value)
}

func (o Option[T]) Unwrap() T {
	if !o.ok {
		panic("option.Unwrap() called on None")
	}

	return o.value
}

func (o Option[T]) Expect(msg string) T {
	if !o.ok {
		panic(msg)
	}

	return o.value
}

func (o Option[T]) UnwrapOr(defaultValue T) T {
	if !o.ok {
		return defaultValue
	}

	return o.value
}

func (o Option[T]) UnwrapOrElse(f func() T) T {
	if !o.ok {
		return f()
	}

	return o.value
}

func (o Option[T]) UnwrapOrZero() T {
	if !o.ok {
		var zero T
		return zero
	}

	return o.value
}

func (o Option[T]) Filter(f func(T) bool) Option[T] {
	if o.ok && f(o.value) {
		return o
	}

	return None[T]()
}

func (o Option[T]) Inspect(f func(T)) Option[T] {
	if o.ok {
		f(o.value)
	}

	return o
}

func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.ok {
		return o
	}

	return other
}

func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.ok {
		return o
	}

	return f()
}

func (o Option[T]) Xor(other Option[T]) Option[T] {
	if o.ok && !other.ok {
		return o
	}

	if !o.ok && other.ok {
		return other
	}

	return None[T]()
}

func (o Option[T]) OkOr(err error) Result[T] {
	if o.ok {
		return Ok(o.value)
	}

	return Err[T](err)
}

func (o Option[T]) OkOrElse(errFn func() error) Result[T] {
	if o.ok {
		return Ok(o.value)
	}

	return Err[T](errFn())
}

func (o Option[T]) ToPtr() *T {
	if !o.ok {
		return nil
	}

	v := o.value
	return &v
}

func (o Option[T]) IfSome(f func(T)) Option[T] {
	if o.ok {
		f(o.value)
	}

	return o
}

func (o Option[T]) IfNone(f func()) Option[T] {
	if !o.ok {
		f()
	}

	return o
}

func (o Option[T]) Match(someFn func(T), noneFn func()) {
	if o.ok {
		someFn(o.value)
	}

	noneFn()
}

func Transpose[T any](o Option[Result[T]]) Result[Option[T]] {
	if !o.ok {
		return Ok(None[T]())
	}

	inner := o.value
	if !inner.ok {
		return Err[Option[T]](inner.err)
	}

	return Ok(Some(inner.value))
}
