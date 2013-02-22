package main

type Flag uint32

const (
  X0 Flag = 1 << iota
  X1
  X2
  X3
  X4
)

type G struct {
  Flags Flag
}

func (f *Flag) Set(flags ...Flag) {
  for _, i := range flags {
    *f = *f | i
  }
}

func (f *Flag) Test(o Flag) bool {
  return (*f & o) != 0
}


func main() {
  println(X0)
  println(X1)
  println(X2)
  println(X3)
  println(X4)
  g := &G{}
  g.Flags.Set(X0)
  g.Flags.Set(X1)
  println("flags:",g.Flags)
  println(g.Flags.Test(X0))
  println(g.Flags.Test(X1))
  println(g.Flags.Test(X2))
  println(g.Flags.Test(X3))
}
