package main

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/Sirupsen/logrus"
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
				Required: true,
			},
		},
	}
}

var res []byte

func resourceAnsibleScriptCreate(d *schema.ResourceData, meta interface{}) error {
	host := d.Get("host").(string)
	runType := d.Get("type").(string)
	file := d.Get("file").(string)

	copyStr := fmt.Sprintf("src=%s dest=/tmp/%s", file, file)
	copyCmd := exec.Command("ansible", host, "-u", "root", "-m", "copy", "-a", copyStr)
	res, err := copyCmd.Output()
	if err != nil {
		logrus.Errorf("error while copy: %s", err)
		return errors.New(string(res))
	}
	logrus.Infof("script copy result: %s", string(res))

	runStr := fmt.Sprintf("%s /tmp/%s", runType, file)
	runCmd := exec.Command("ansible", host, "-u", "root", "-a", runStr)
	res, err = runCmd.Output()
	if err != nil {
		logrus.Errorf("error while execute: %s", err)
		return errors.New(string(res))
	}
	logrus.Infof("script run result: %s", string(res))
	d.Set("result", string(res))
	d.SetId("1")

	return resourceAnsibleScriptRead(d, meta)

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
