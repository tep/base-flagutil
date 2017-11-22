package flagutil

func (g *FlagsGroup) SetUnless(name, value string) (bool, error) {
	if _, ok := g.ValueIsSet(name); ok {
		return false, nil
	}

	if err := g.Set(name, value); err != nil {
		return false, err
	}

	return true, nil
}
