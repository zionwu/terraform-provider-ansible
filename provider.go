package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_URL", ""),
				Description: descriptions["api_url"],
			},
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			},
			"config": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_CLIENT_CONFIG", ""),
				Description: descriptions["config"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ansible_script": resourceAnsibleScript(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return nil, nil
}
