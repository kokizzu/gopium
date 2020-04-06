package strategies

import (
	"context"

	"1pkg/gopium"
)

// list of nope presets
var (
	np = nope{}
)

// nope defines nil strategy implementation
// that does nothing by returning original structure
type nope struct{}

// Apply nope implementation
func (stg nope) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	return o, nil
}