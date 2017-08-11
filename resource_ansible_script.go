package main

import (
	"fmt"
	"os/exec"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAnsibleScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceAnsibleScriptCreate,
		Read:   resourceAnsibleScriptRead,
		Update: resourceAnsibleScriptUpdate,
		Delete: resourceAnsibleScriptDelete,

		Schema: map[string]*schema.Schema{
			"file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"result": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAnsibleScriptCreate(d *schema.ResourceData, meta interface{}) error {
	host := d.Get("host").(string)
	runType := d.Get("type").(string)
	file := d.Get("file").(string)

	copyStr := fmt.Sprintf("src=%s dest=/tmp/%s", file, file)
	copyCmd := exec.Command("ansible", host, "-u", "root", "-m", "copy", "-a", copyStr)
	res, err := copyCmd.Output()
	if err != nil {
		return err
	}

	runStr := fmt.Sprintf("%s /tmp/%s", runType, file)
	runCmd := exec.Command("ansible", host, "-u", "root", "-a", runStr)
	res, err = runCmd.Output()
	if err != nil {
		return err
	}

	d.Set("result", res)
	return nil

}

func resourceAnsibleScriptRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAnsibleScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAnsibleScriptDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
