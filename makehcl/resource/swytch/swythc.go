package swytch

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/x86taka/sacloud-tools/makehcl/resource/utils"
	"github.com/zclconf/go-cty/cty"
)

type LocalSwitchHCL struct {
	ResourceName    string
	Name            string
	Tags            []string
	Size            int64
	SourceArchiveID string
}

func (v *LocalSwitchHCL) OutputHCL() string {
	f := hclwrite.NewEmptyFile()

	rootBody := f.Body()
	moduleBlock := rootBody.AppendNewBlock("resource", []string{"sakuracloud_switch", utils.FormatHCL(v.ResourceName)})
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

	// TimeOuts
	moduleBody.AppendNewline()
	timeouts := moduleBody.AppendNewBlock("timeouts", []string{})
	timeoutsBody := timeouts.Body()
	timeoutsBody.SetAttributeValue("create", cty.StringVal("1h"))
	timeoutsBody.SetAttributeValue("delete", cty.StringVal("1h"))
	return string(f.Bytes())
}
