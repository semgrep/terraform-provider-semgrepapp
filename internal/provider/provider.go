package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SEMGREP_APP_TOKEN", nil),
			},
		},
		p := &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"semgrepapp_rule":    dataSourceRules(),
			"semgrepapp_ruleset": dataSourceRulesets(),
		},
			ResourcesMap: map[string]*schema.Resource{
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	isAuthenticated bool
	token           string
	userAgent			  string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		token := d.Get("token").(string)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if token != "" {
			return &apiClient{true, token, userAgent}, diags
		}

		return &apiClient{false, "", userAgent}, diags
	}
}
