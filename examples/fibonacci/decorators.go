package fibonacci

// Decorator decorates a Provider.
type Decorator func(Provider) Provider

// Memoize decorates p with memoization.
var Memoize Decorator = func(p Provider) Provider {
	return NewMemoized(p)
}

// NoDecorator is a noop Decorator.
var NoDecorator Decorator = func(p Provider) Provider {
	return p
}
