package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"time"

	"github.com/pkg/errors"

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
			"source_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"target_path": &schema.Schema{
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
			"param": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sleep_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1200,
			},
			"result": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

var res []byte

func resourceAnsibleScriptCreate(d *schema.ResourceData, meta interface{}) error {
	host := d.Get("host").(string)
	hostUsername := d.Get("host_username").(string)
	hostPassword := d.Get("host_password").(string)
	runType := d.Get("type").(string)
	targetPath := d.Get("target_path").(string)
	sourcePath := d.Get("source_path").(string)
	param := d.Get("param").(string)
	sleepInterval := d.Get("sleep_interval").(int)

	dial("tcp", host+":22", time.Duration(sleepInterval)*time.Second)

	//write ansible host config
	f, err := os.OpenFile("/etc/ansible/hosts", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		logrus.Errorf("error while copy: %s", err)
		d.SetId("1")
		return err
	}
	defer f.Close()

	config := fmt.Sprintf("%s ansible_ssh_user=%s ansible_ssh_pass=%s", host, hostUsername, hostPassword)
	hostExist := false
	scanner := bufio.NewScanner(f)
	logrus.Info(config)
	for scanner.Scan() {
		line := scanner.Text()

		logrus.Info(line)
		if line == config {
			hostExist = true
			break
		}
	}

	if !hostExist {
		_, err := f.WriteString(config + "\n")
		if err != nil {
			logrus.Errorf("error while updating ansible host: %s", err)
			d.SetId("1")
			return err
		}
	}
	file := filepath.Base(sourcePath)
	copyStr := fmt.Sprintf("src=%s dest=%s", sourcePath, filepath.Join(targetPath, file))
	copyCmd := exec.Command("ansible", host, "-u", hostUsername, "-m", "copy", "-a", copyStr)
	resCopy, err := copyCmd.Output()
	logrus.Infof("exec %v", copyCmd.Args)
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			err = errors.Wrapf(err, "%s", ee.Stderr)
		}
		logrus.Errorf("error while copy: %s res: %s", err, string(resCopy))
		d.Set("result", string(resCopy))
		d.SetId("1")
		return err
	}
	logrus.Infof("script copy result: %s", string(resCopy))

	runStr := fmt.Sprintf("%s %s %s", runType, filepath.Join(targetPath, file), param)
	runCmd := exec.Command("ansible", host, "-u", hostUsername, "-a", runStr)
	res, err := runCmd.Output()
	logrus.Infof("exec %v", runCmd.Args)

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			err = errors.Wrapf(err, "%s", ee.Stderr)
		}
		logrus.Errorf("error while execute: %s res: %s", err, string(res))
		d.Set("result", string(res))
		d.SetId("1")
		return err
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

// Dial dial the raddr before timeout
func dial(protocol, raddr string, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			return false
		default:
			conn, err := net.DialTimeout(protocol, raddr, time.Second)
			if err == nil {
				conn.Close()
				return true
			}
			time.Sleep(time.Second)
		}
	}
}
