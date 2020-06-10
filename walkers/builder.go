package walkers

import (
	"fmt"

	"1pkg/gopium/gopium"
)

// list of registered types walkers
const (
	// wast walkers
	AstStd    gopium.WalkerName = "ast_std"
	AstGo     gopium.WalkerName = "ast_go"
	AstGoTree gopium.WalkerName = "ast_go_tree"
	AstGopium gopium.WalkerName = "ast_gopium"
	// wout walkers
	JsonbFile gopium.WalkerName = "json_file"
	XmlbFile  gopium.WalkerName = "xml_file"
	CsvbFile  gopium.WalkerName = "csv_file"
	MdtFile   gopium.WalkerName = "md_table_file"
	// wdiff walkers
	SizeAlignMdtFile gopium.WalkerName = "size_align_md_table_file"
	FieldsHtmltFile  gopium.WalkerName = "fields_html_table_file"
)

// Builder defines types gopium.WalkerBuilder implementation
// that uses parser and exposer to pass it to related walkers
type Builder struct {
	Parser  gopium.Parser  `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	Exposer gopium.Exposer `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	Printer gopium.Printer `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	Deep    bool           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	Bref    bool           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	_       [14]byte       `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; - 🌺 gopium @1pkg

// Build Builder implementation
func (b Builder) Build(name gopium.WalkerName) (gopium.Walker, error) {
	switch name {
	// wast walkers
	case AstStd:
		return aststd.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	case AstGo:
		return astgo.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	case AstGoTree:
		return astgotree.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	case AstGopium:
		return astgopium.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	// wout walkers
	case JsonbFile:
		return jsonfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case XmlbFile:
		return xmlfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case CsvbFile:
		return csvfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case MdtFile:
		return mdtfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	// wdiff walkers
	case SizeAlignMdtFile:
		return samdtfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case FieldsHtmltFile:
		return fhtmltfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	default:
		return nil, fmt.Errorf("walker %q wasn't found", name)
	}
}
