package markup

// If returns nil if cond is false, otherwise it returns the given children.
func If(cond bool, children ...MarkupOrChild) MarkupOrChild {
	if len(children) == 0 {
		return nil
	}
	if cond {
		if len(children) == 1 {
			return children[0]
		}
		return List(children)
	}
	return nil
}
