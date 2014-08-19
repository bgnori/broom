package ps

func IterOverPairAsList(v Value) chan Value {
	ch := make(chan Value)
	go func() {
		defer close(ch)
		for v != nil && v.Pair() { // ignores last cdr value.
			pair, _ := v.(*Pair)
			ch <- pair.Car()
			v = pair.Cdr()
		}
	}()
	return ch
}

type Z struct {
	x, y Value
}

func Zip(xs, ys Value) chan Z {
	ch := make(chan Z)
	go func() {
		defer close(ch)
		for xs != nil && ys != nil && xs.Pair() && ys.Pair() {
			x, _ := xs.(*Pair)
			y, _ := ys.(*Pair)
			ch <- Z{x: x.Car(), y: y.Car()}
			xs = x.Cdr()
			ys = y.Cdr()
		}
	}()
	return ch
}

func RecEq(x, y Value) bool {
    if x == nil && y == nil {
        return true
    }
    if x == nil || y == nil {
        return false
    }
    if x.Pair() && y.Pair() {
        p, _ := x.(*Pair)
        q, _ := y.(*Pair)
        return RecEq(p.Car(), q.Car()) && RecEq(p.Cdr(), q.Cdr())
    }
    if x.Pair() || y.Pair() {
        return false
    }
    return x == y
}

func ch2xs(ch chan *Token) []*Token {
    xs := make([]*Token, 0)
    for t := range ch {
        xs = append(xs, t)
    }
    return xs
}


