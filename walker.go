package gopium

import (
	"context"
	"regexp"
)

// Walker defines hierarchical walker abstraction
// that applies some strategy to tree structures
// and modifies tree or creates other side effects
type Walker interface {
	// VisitTop visits only top level of the tree
	VisitTop(context.Context, *regexp.Regexp, Strategy) error
	// VisitDeep visits all levels of the tree
	VisitDeep(context.Context, *regexp.Regexp, Strategy) error
}

// WalkerName defines registred Walker name abstraction
// used by WalkerBuilder to build registred walkers
type WalkerName string

// WalkerBuilder defines Walker builder abstraction
// that helps to create Walker by WalkerName
type WalkerBuilder interface {
	Build(StrategyName) (Walker, error)
}