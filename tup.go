package wire

import "fmt"

type Tup2Val[T0, T1 any] struct {
	E0 T0
	E1 T1
}

func Tup2[T0, T1 any](p0 Proto[T0], p1 Proto[T1]) Proto[Tup2Val[T0, T1]] {
	return proto[Tup2Val[T0, T1]]{
		read:  readTup2(p0, p1),
		write: writeTup2(p0, p1),
		size:  sizeTup2(p0, p1),
	}
}

func readTup2[T0, T1 any](p0 Proto[T0], p1 Proto[T1]) func(Reader) (Tup2Val[T0, T1], error) {
	return func(r Reader) (Tup2Val[T0, T1], error) {
		e0, err := p0.Read(r)
		if err != nil {
			return Tup2Val[T0, T1]{}, err
		}
		e1, err := p1.Read(r)
		if err != nil {
			return Tup2Val[T0, T1]{}, err
		}
		return Tup2Val[T0, T1]{e0, e1}, nil
	}
}

func writeTup2[T0, T1 any](p0 Proto[T0], p1 Proto[T1]) func(Writer, Tup2Val[T0, T1]) error {
	return func(w Writer, tup Tup2Val[T0, T1]) error {
		if err := p0.Write(w, tup.E0); err != nil {
			return fmt.Errorf("Tup2.Write: E0: %w", err)
		}
		if err := p1.Write(w, tup.E1); err != nil {
			return fmt.Errorf("Tup2.Write: E1: %w", err)
		}
		return nil
	}
}

func sizeTup2[T0, T1 any](p0 Proto[T0], p1 Proto[T1]) func(Tup2Val[T0, T1]) uint64 {
	return func(tup Tup2Val[T0, T1]) uint64 { return p0.Size(tup.E0) + p1.Size(tup.E1) }
}

type Tup3Val[T0, T1, T2 any] struct {
	E0 T0
	E1 T1
	E2 T2
}

func Tup3[T0, T1, T2 any](p0 Proto[T0], p1 Proto[T1], p2 Proto[T2]) Proto[Tup3Val[T0, T1, T2]] {
	return proto[Tup3Val[T0, T1, T2]]{
		read:  readTup3(p0, p1, p2),
		write: writeTup3(p0, p1, p2),
		size:  func(tup Tup3Val[T0, T1, T2]) uint64 { return p0.Size(tup.E0) + p1.Size(tup.E1) + p2.Size(tup.E2) },
	}
}

func readTup3[T0, T1, T2 any](p0 Proto[T0], p1 Proto[T1], p2 Proto[T2]) func(Reader) (Tup3Val[T0, T1, T2], error) {
	return func(r Reader) (Tup3Val[T0, T1, T2], error) {
		e0, err := p0.Read(r)
		if err != nil {
			return Tup3Val[T0, T1, T2]{}, err
		}
		e1, err := p1.Read(r)
		if err != nil {
			return Tup3Val[T0, T1, T2]{}, err
		}
		e2, err := p2.Read(r)
		if err != nil {
			return Tup3Val[T0, T1, T2]{}, err
		}
		return Tup3Val[T0, T1, T2]{e0, e1, e2}, nil
	}
}

func writeTup3[T0, T1, T2 any](p0 Proto[T0], p1 Proto[T1], p2 Proto[T2]) func(Writer, Tup3Val[T0, T1, T2]) error {
	return func(w Writer, tup Tup3Val[T0, T1, T2]) error {
		if err := p0.Write(w, tup.E0); err != nil {
			return fmt.Errorf("Tup3.Write: E0: %w", err)
		}
		if err := p1.Write(w, tup.E1); err != nil {
			return fmt.Errorf("Tup3.Write: E1: %w", err)
		}
		if err := p2.Write(w, tup.E2); err != nil {
			return fmt.Errorf("Tup3.Write: E2: %w", err)
		}
		return nil
	}
}
