package expt

func WithSelectors(selectors ...Selector) func(e *Expt) {
	return func(e *Expt) {
		e.selectors = append(e.selectors, selectors...)
	}
}

func WithFilters(filters ...ExptFilter) func(e *Expt) {
	return func(e *Expt) {
		e.filters = append(e.filters, filters...)
	}
}

func WithTrafficer(trafficer Trafficer) func(e *Expt) {
	return func(e *Expt) {
		e.trafficer = trafficer
	}
}
