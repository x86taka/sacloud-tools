package vm

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/x86taka/sacloud-tools/makehcl/resource/utils"
	"github.com/zclconf/go-cty/cty"
)

type VM struct {
	ResourceName string
	Name         string
	Tags         []string
	Cpus         int64
	Mem          int64
	Disks        map[int]string
	Nics         map[int]string
}

func (v *VM) OutputHCL() string {
	f := hclwrite.NewEmptyFile()

	rootBody := f.Body()
	moduleBlock := rootBody.AppendNewBlock("resource", []string{"sakuracloud_server", utils.FormatHCL(v.ResourceName)})
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
	moduleBody.SetAttributeValue("core", cty.NumberIntVal(v.Cpus))
	moduleBody.SetAttributeValue("memory", cty.NumberIntVal(v.Mem))
	moduleBody.AppendNewline()
	// Disk
	diskList := "["
	for i := 0; i < len(v.Disks); i++ {
		diskList += v.Disks[0]
	}
	diskList += "]"
	moduleBody.SetAttributeRaw("disks", hclwrite.Tokens{&hclwrite.Token{Bytes: []byte(diskList)}})
	// Nic
	moduleBody.AppendNewline()
	for i := 0; i < len(v.Nics); i++ {
		nic := moduleBody.AppendNewBlock("network_interface", []string{})
		nic.Body().SetAttributeRaw("upstream", hclwrite.Tokens{&hclwrite.Token{
			Bytes: []byte(v.Nics[i]),
		}})
	}
	// TimeOuts
	moduleBody.AppendNewline()
	timeouts := moduleBody.AppendNewBlock("timeouts", []string{})
	timeoutsBody := timeouts.Body()
	timeoutsBody.SetAttributeValue("create", cty.StringVal("1h"))
	timeoutsBody.SetAttributeValue("delete", cty.StringVal("1h"))
	return string(f.Bytes())
}
