package godifft

type Change int

const (
	Insert Change = iota
	Remove
	Keep
)

type Edit[T any] struct {
	Change  Change
	Element T
}

type DiffTOptions[T any] struct {
	Equals func(T, T) bool
}

func DiffT[T any](xs, ys []T, opts DiffTOptions[T]) []Edit[T] {
	eq := opts.Equals
	if eq == nil {
		eq = func(x T, y T) bool {
			return any(x) == any(y)
		}
	}
	d := &differ[T]{opts.Equals, xs, ys}
	return d.diff()
}

type differ[T any] struct {
	eq func(T, T) bool
	xs []T
	ys []T
}

func (d *differ[T]) difflen() *matrix {
	difflen := newMatrix(len(d.xs)+1, len(d.ys)+1)
	for xp := len(d.xs); xp >= 0; xp-- {
		for yp := len(d.ys); yp >= 0; yp-- {
			l, _ := d.choose(difflen, xp, yp)
			difflen.set(xp, yp, l)
		}
	}
	return difflen
}

func (d *differ[T]) choose(difflen *matrix, xp, yp int) (int, Change) {
	xrem := len(d.xs) - xp
	yrem := len(d.ys) - yp
	switch {
	case xrem == 0:
		return yrem, Insert
	case yrem == 0:
		return xrem, Remove
	}
	l := 1 + difflen.get(xp+1, yp)
	c := Remove
	if n := 1 + difflen.get(xp, yp+1); n < l {
		l = n
		c = Insert
	}
	if d.eq(d.xs[xp], d.ys[yp]) {
		if n := difflen.get(xp+1, yp+1); n < l {
			l = n
			c = Keep
		}
	}
	return l, c
}

func (d *differ[T]) diff() []Edit[T] {
	var edits []Edit[T]
	difflen, xs, ys := d.difflen(), d.xs, d.ys
	for {
		if len(xs) == 0 {
			for _, y := range ys {
				edits = append(edits, d.insert(y))
			}
			return edits
		}
		if len(ys) == 0 {
			for _, x := range xs {
				edits = append(edits, d.remove(x))
			}
			return edits
		}
		xp, yp := len(d.xs)-len(xs), len(d.ys)-len(ys)
		_, diff := d.choose(difflen, xp, yp)
		switch diff {
		case Remove:
			edits, xs = append(edits, d.remove(xs[0])), xs[1:]
		case Insert:
			edits, ys = append(edits, d.insert(ys[0])), ys[1:]
		default: // keep
			edits, xs, ys = append(edits, d.keep(xs[0])), xs[1:], ys[1:]
		}
	}
}

func (d *differ[T]) insert(x T) Edit[T] {
	return Edit[T]{Insert, x}
}

func (d *differ[T]) remove(x T) Edit[T] {
	return Edit[T]{Remove, x}
}

func (d *differ[T]) keep(x T) Edit[T] {
	return Edit[T]{Keep, x}
}

type Differ[T any, Diff any] interface {
	Added(T) Diff
	Removed(T) Diff
	Diff(T, T) (Diff, bool)
}

func DiffMapT[K comparable, T, R any](
	differ Differ[T, R],
	m1, m2 map[K]T,
) map[K]R {
	diffmap := map[K]R{}
	for k, xvv := range m1 {
		yvv, ok := m2[k]
		if !ok {
			diffmap[k] = differ.Removed(xvv)
		} else {
			xdy, xneqy := differ.Diff(xvv, yvv)
			if xneqy {
				diffmap[k] = xdy
			}
		}
	}
	for k, yvv := range m2 {
		if _, ok := m1[k]; !ok {
			diffmap[k] = differ.Added(yvv)
		}
	}
	return diffmap
}

type treeDiffer struct {
	innerDiffer Differ[interface{}, interface{}]
	equals      func(interface{}, interface{}) bool
}

var _ Differ[interface{}, interface{}] = &treeDiffer{}

func (t *treeDiffer) Added(x interface{}) interface{} {
	return t.innerDiffer.Added(x)
}

func (t *treeDiffer) Removed(x interface{}) interface{} {
	return t.innerDiffer.Removed(x)
}

func (t *treeDiffer) Diff(tree1, tree2 interface{}) (interface{}, bool) {
	switch xv := tree1.(type) {
	case []interface{}:
		switch yv := tree2.(type) {
		case []interface{}:
			d := []interface{}{}
			diffs := DiffT(xv, yv, DiffTOptions[interface{}]{t.equals})
			eq := true
			for _, ed := range diffs {
				switch ed.Change {
				case Insert:
					d = append(d, t.Added(ed.Element))
					eq = false
				case Remove:
					d = append(d, t.Removed(ed.Element))
					eq = false
				case Keep:
					d = append(d, ed.Element)
				}
			}
			return d, !eq
		default:
			return t.innerDiffer.Diff(xv, yv)
		}
	case map[string]interface{}:
		switch yv := tree2.(type) {
		case map[string]interface{}:
			m := DiffMapT(t, xv, yv)
			return m, len(m) > 0
		default:
			return t.innerDiffer.Diff(tree1, tree2)
		}
	default:
		return t.innerDiffer.Diff(tree1, tree2)
	}
}

func DiffTree[Diff any](
	differ Differ[interface{}, interface{}],
	equals func(interface{}, interface{}) bool,
	tree1, tree2 interface{},
) (interface{}, bool) {
	td := &treeDiffer{differ, equals}
	return td.Diff(tree1, tree2)
}
