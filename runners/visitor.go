package runners

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"1pkg/gopium"
)

// visitor defines helper
// that coordinates runner stages
// to finaly make visiting
// - strategy building
// - walker building
// - visiting
type visitor struct {
	regex   *regexp.Regexp
	timeout time.Duration
}

// strategy builds strategy instance
// by using builder and strategies names
func (visitor) strategy(b gopium.StrategyBuilder, snames []gopium.StrategyName) (gopium.Strategy, error) {
	// build strategy
	stg, err := b.Build(snames...)
	if err != nil {
		return nil, fmt.Errorf("can't build such strategy %v %v", snames, err)
	}
	return stg, nil
}

// walker builds walker instance
// by using builder and walker name
func (visitor) walker(b gopium.WalkerBuilder, wname gopium.WalkerName) (gopium.Walker, error) {
	// build walker
	walker, err := b.Build(wname)
	if err != nil {
		return nil, fmt.Errorf("can't build such walker %q %v", wname, err)
	}
	return walker, nil
}

// visit coordinates walker visiting
func (v visitor) visit(ctx context.Context, w gopium.Walker, stg gopium.Strategy) error {
	// set up timeout context
	if v.timeout > 0 {
		nctx, cancel := context.WithTimeout(ctx, v.timeout)
		defer cancel()
		ctx = nctx
	}
	// exec visit on walker with strategy
	if err := w.Visit(ctx, v.regex, stg); err != nil {
		return fmt.Errorf("visiting error happened %v", err)
	}
	return nil
}