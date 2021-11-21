package semgrepapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRulesRead,
		Schema: map[string]*schema.Schema{
			"rules": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"visibility": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha_sum": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_uri": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"deployment_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"meta": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*AppContext)
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/registry/rule", "https://semgrep.dev/api"), nil)
	if c.isAuthenticated {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	rules := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&rules)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unexpected response from Semgrep App",
			Detail:   fmt.Sprintf("The semgrep.dev API returned an unexpected response with %s status.", r.Status),
		})

		return diags
	}

	if err := d.Set("rules", rules); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
