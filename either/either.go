package either

import "errors"

type Either[L, R any] struct {
	left  L
	right R

	isLeft bool
}

func Left[L, R any](value L) Either[L, R] {
	return Either[L, R]{left: value, isLeft: true}
}

func Right[L, R any](value R) Either[L, R] {
	return Either[L, R]{right: value, isLeft: false}
}

func (e Either[L, R]) IsLeft() bool { return e.isLeft }

func (e Either[L, R]) IsRight() bool { return !e.isLeft }

func (e Either[L, R]) IsLeftAnd(f func(L) bool) bool {
	return e.isLeft && f(e.left)
}

func (e Either[L, R]) IsRightAnd(f func(R) bool) bool {
	return !e.isLeft && f(e.right)
}

func (e Either[L, R]) UnwrapLeft() L {
	if !e.isLeft {
		panic("either.UnwrapLeft: called on Right")
	}

	return e.left
}

func (e Either[L, R]) UnwrapRight() R {
	if e.isLeft {
		panic("either.UnwrapRight: called on Left")
	}

	return e.right
}

func (e Either[L, R]) UnwrapLeftOr(defaultValue L) L {
	if !e.isLeft {
		return defaultValue
	}

	return e.left
}

func (e Either[L, R]) UnwrapRightOr(defaultValue R) R {
	if e.isLeft {
		return defaultValue
	}

	return e.right
}

func (e Either[L, R]) UnwrapLeftOrElse(f func(R) L) L {
	if !e.isLeft {
		return f(e.right)
	}

	return e.left
}

func (e Either[L, R]) UnwrapRightOrElse(f func(L) R) R {
	if e.isLeft {
		return f(e.left)
	}

	return e.right
}

func MapLeft[L, R, L2 any](e Either[L, R], f func(L) L2) Either[L2, R] {
	if e.isLeft {
		return Left[L2, R](f(e.left))
	}

	return Right[L2](e.right)
}

func MapRight[L, R, R2 any](e Either[L, R], f func(R) R2) Either[L, R2] {
	if !e.isLeft {
		return Right[L](f(e.right))
	}

	return Left[L, R2](e.left)
}

func MapBoth[L, R, L2, R2 any](e Either[L, R], leftFn func(L) L2, rightFn func(R) R2) Either[L2, R2] {
	if e.isLeft {
		return Left[L2, R2](leftFn(e.left))
	}

	return Right[L2](rightFn(e.right))
}

func FlatMapLeft[L, R, L2 any](e Either[L, R], f func(L) Either[L2, R]) Either[L2, R] {
	if e.isLeft {
		return f(e.left)
	}

	return Right[L2](e.right)
}

func FlatMapRight[L, R, R2 any](e Either[L, R], f func(R) Either[L, R2]) Either[L, R2] {
	if !e.isLeft {
		return f(e.right)
	}

	return Left[L, R2](e.left)
}

func Fold[L, R, T any](e Either[L, R], leftFn func(L) T, rightFn func(R) T) T {
	if e.isLeft {
		return leftFn(e.left)
	}

	return rightFn(e.right)
}

func (e Either[L, R]) Swap() Either[R, L] {
	if e.isLeft {
		return Right[R](e.left)
	}

	return Left[R, L](e.right)
}

func (e Either[L, R]) IfLeft(f func(L)) Either[L, R] {
	if e.isLeft {
		f(e.left)
	}

	return e
}

func (e Either[L, R]) IfRight(f func(R)) Either[L, R] {
	if !e.isLeft {
		f(e.right)
	}

	return e
}

func (e Either[L, R]) Tap(leftFn func(L), rightFn func(R)) Either[L, R] {
	if e.isLeft {
		leftFn(e.left)
	} else {
		rightFn(e.right)
	}

	return e
}

func ContainsLeft[L comparable, R any](e Either[L, R], target L) bool {
	return e.isLeft && e.left == target
}

func ContainsRight[L any, R comparable](e Either[L, R], target R) bool {
	return !e.isLeft && e.right == target
}

func FromResult[T any, E interface{ Error() string }](isOk bool, value T, err error) Either[T, error] {
	if err == nil {
		return Left[T, error](value)
	}

	return Right[T](err)
}

func ToResult[T any](e Either[T, error]) (T, error) {
	if e.isLeft {
		return e.left, nil
	}

	var zero T
	return zero, e.right
}

func LeftOr[T any](e Either[T, error]) (T, error) {
	return ToResult(e)
}

func Partition[L, R any](eithers []Either[L, R]) (lefts []L, rights []R) {
	for _, e := range eithers {
		if e.isLeft {
			lefts = append(lefts, e.left)
		} else {
			rights = append(rights, e.right)
		}
	}

	return lefts, rights
}

func Merge[T any](e Either[T, T]) T {
	if e.isLeft {
		return e.left
	}

	return e.right
}

func ErrRight[L any](e Either[L, error]) error {
	if !e.isLeft {
		return e.right
	}

	return nil
}

func IsErr[L any](e Either[L, error], target error) bool {
	if e.isLeft {
		return false
	}

	return errors.Is(e.right, target)
}
