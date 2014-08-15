package ps

func IterOverPairAsList(v Value) chan Value {
        ch := make(chan Value)
        go func() {
                  defer close(ch)
                  for v!= nil && v.Pair() { // ignores last cdr value.
                          pair, _ := v.(*Pair)
                          ch<- pair.Car()
                          v = pair.Cdr()
                  }
        }()
        return ch
}


type Z struct{
  x, y Value
}

func Zip(xs, ys Value) chan Z {
        ch := make(chan Z)
        go func() {
                  defer close(ch)
                  for xs !=nil && ys != nil && xs.Pair() && ys.Pair() {
                          x, _ := xs.(*Pair)
                          y, _ := ys.(*Pair)
                          ch<- Z{x:x.Car(), y:y.Car()}
                          xs = x.Cdr()
                          ys = y.Cdr()
                  }
        }()
        return ch
}


