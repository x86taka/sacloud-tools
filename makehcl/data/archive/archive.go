package archive

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/x86taka/sacloud-tools/makehcl/resource/utils"
	"github.com/zclconf/go-cty/cty"
)

type DataArchiveHCL struct {
	Name string
}

func (v *DataArchiveHCL) OutputHCL() string {
	f := hclwrite.NewEmptyFile()

	rootBody := f.Body()
	moduleBlock := rootBody.AppendNewBlock("data", []string{"sakuracloud_archive", utils.FormatHCL(v.Name)})
	moduleBody := moduleBlock.Body()
	filter := moduleBody.AppendNewBlock("filter", []string{})
	filter.Body().SetAttributeValue("names", cty.ListVal([]cty.Value{cty.StringVal(v.Name)}))

	return string(f.Bytes())
}
