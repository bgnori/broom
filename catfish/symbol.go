package catfish

type symbolImpl struct {
	value string
}

func sym(s string) *symbolImpl {
	return &symbolImpl{value: s}
}

func (s *symbolImpl) GetValue() string {
	return s.value
}

func (s *symbolImpl) Eq(other Value) bool {
	if t, ok := other.(*symbolImpl); ok {
		return s.value == t.value
	} else {
		return false
	}
}

func (s *symbolImpl) String () string {
      return s.value
}
