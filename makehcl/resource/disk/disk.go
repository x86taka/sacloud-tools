package disk

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/x86taka/sacloud-tools/makehcl/resource/utils"
	"github.com/zclconf/go-cty/cty"
)

type DiskHCL struct {
	Name            string
	Tags            []string
	Size            int64
	SourceArchiveID string
}

func (v *DiskHCL) OutputHCL() string {
	f := hclwrite.NewEmptyFile()

	rootBody := f.Body()
	moduleBlock := rootBody.AppendNewBlock("resource", []string{"sakuracloud_disk", utils.FormatHCL(v.Name)})
	moduleBody := moduleBlock.Body()

	nameToken := hclwrite.Tokens{
		{
			Bytes: []byte("\"" + v.Name + "\""),
		},
	}

	moduleBody.SetAttributeRaw("name", nameToken)

	tagList := hclwrite.Tokens{
		{
			Bytes: []byte("["),
		},
	}
	for _, tag := range v.Tags {
		token := &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("\"" + tag + "\","),
		}
		tagList = append(tagList, token)
	}
	tagList = append(tagList, &hclwrite.Token{
		Bytes: []byte("]"),
	})

	moduleBody.SetAttributeRaw("tags", tagList)
	moduleBody.AppendNewline()
	sourceArchiveToken := hclwrite.Tokens{
		{
			Type: hclsyntax.TokenIdent,
			//Bytes: []byte("\"${var.team_id}-${local.problem_name}-" + v.Name + "\""),
			Bytes: []byte(v.SourceArchiveID),
		},
	}
	moduleBody.SetAttributeRaw("source_archive_id", sourceArchiveToken)
	moduleBody.SetAttributeValue("size", cty.NumberIntVal(v.Size))

	// TimeOuts
	moduleBody.AppendNewline()
	timeouts := moduleBody.AppendNewBlock("timeouts", []string{})
	timeoutsBody := timeouts.Body()
	timeoutsBody.SetAttributeValue("create", cty.StringVal("1h"))
	timeoutsBody.SetAttributeValue("delete", cty.StringVal("1h"))
	return string(f.Bytes())
}
