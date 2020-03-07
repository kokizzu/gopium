package strategy

import (
	"context"
	"go/types"

	"1pkg/gopium"
)

// enum defines struct enumerating strategy implementation
// that goes through all structure fields and uses gopium.Whistleblower
// to expose gopium.Field DTO for each field
// and puts it back to resulted gopium.Struct object
type enum struct {
	wb gopium.Whistleblower
}

// Apply enum implementation
func (stg enum) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// set structure name
	r.Name = name
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Fields = make([]gopium.Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// get tag
		tag := st.Tag(i)
		// get typeinfo
		tname, tsize := stg.wb.Expose(f.Type())
		// fill field structure
		r.Fields = append(r.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     tname,
			Size:     tsize,
			Tag:      tag,
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	o = r
	return
}
