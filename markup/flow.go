package markup

// If returns nil if cond is false, otherwise it returns the given children.
func If(cond bool, children ...MarkupOrChild) MarkupOrChild {
	if cond {
		return List(children)
	}
	return nil
}
