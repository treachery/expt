package expt

func WithSelectors(selectors ...Selector) func(e *Expt) {
	return func(e *Expt) {
		e.Selectors = append(e.Selectors, selectors...)
	}
}

func WithFilters(filters ...ExptFilter) func(e *Expt) {
	return func(e *Expt) {
		e.Filters = append(e.Filters, filters...)
	}
}
