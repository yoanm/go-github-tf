package core

import (
	"fmt"
	"testing"

	differ "github.com/andreyvit/diff"
	"github.com/google/go-cmp/cmp"
)

func TestLoadTemplateList(t *testing.T) {
	TemplateLevel0 := "a-template0"
	TemplateLevel1 := "a-template1"
	TemplateLevel2 := "a-template2"
	TemplateLevel3 := "a-template3"
	TemplateLevel4 := "a-template4"
	TemplateLevel5 := "a-template5"
	TemplateLevel6 := "a-template6"
	TemplateLevel7 := "a-template7"
	TemplateLevel8 := "a-template8"
	TemplateLevel9 := "a-template9"
	TemplateLevel10 := "a-template10"
	templateLevel11 := "a-template11"
	a11SubTemplatesTemplate := "a-template12"
	description0 := "my description0"
	description1 := "my description1"
	description2 := "my description2"
	description3 := "my description3"
	description4 := "my description4"
	description5 := "my description5"
	description6 := "my description6"
	description7 := "my description7"
	description8 := "my description8"
	description9 := "my description9"
	description10 := "my description10"
	description11 := "my description11"
	// emptyTplConfig := &TemplatesConfig{}
	tplConfig := &TemplatesConfig{
		Repos: map[string]*GhRepoConfig{
			TemplateLevel0:  {Description: &description0},
			TemplateLevel1:  {ConfigTemplates: &[]string{TemplateLevel0}, Description: &description1},
			TemplateLevel2:  {ConfigTemplates: &[]string{TemplateLevel1}, Description: &description2},
			TemplateLevel3:  {ConfigTemplates: &[]string{TemplateLevel2}, Description: &description3},
			TemplateLevel4:  {ConfigTemplates: &[]string{TemplateLevel3}, Description: &description4},
			TemplateLevel5:  {ConfigTemplates: &[]string{TemplateLevel4}, Description: &description5},
			TemplateLevel6:  {ConfigTemplates: &[]string{TemplateLevel5}, Description: &description6},
			TemplateLevel7:  {ConfigTemplates: &[]string{TemplateLevel6}, Description: &description7},
			TemplateLevel8:  {ConfigTemplates: &[]string{TemplateLevel7}, Description: &description8},
			TemplateLevel9:  {ConfigTemplates: &[]string{TemplateLevel8}, Description: &description9},
			TemplateLevel10: {ConfigTemplates: &[]string{TemplateLevel9}, Description: &description10},
			templateLevel11: {ConfigTemplates: &[]string{TemplateLevel10}, Description: &description11},
			a11SubTemplatesTemplate: {
				ConfigTemplates: &[]string{
					TemplateLevel1, TemplateLevel2, TemplateLevel3, TemplateLevel4, TemplateLevel5,
					TemplateLevel6, TemplateLevel7, TemplateLevel8, TemplateLevel9, TemplateLevel10,
					templateLevel11,
				},
			},
		},
	}
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"more than 10 templates to load": {
			&GhRepoConfig{
				ConfigTemplates: &[]string{
					TemplateLevel1, TemplateLevel2, TemplateLevel3, TemplateLevel4, TemplateLevel5,
					TemplateLevel6, TemplateLevel7, TemplateLevel8, TemplateLevel9, TemplateLevel10,
					templateLevel11,
				},
			},
			tplConfig,
			nil,
			fmt.Errorf("more than 10 repository template detected for ROOT"),
		},
		"leaf with more than 10 templates to load": {
			&GhRepoConfig{
				ConfigTemplates: &[]string{TemplateLevel1, a11SubTemplatesTemplate, TemplateLevel2},
			},
			tplConfig,
			nil,
			fmt.Errorf("more than 10 repository template detected for a-template12"),
		},
		"more than 10 level of template to load": {
			&GhRepoConfig{
				ConfigTemplates: &[]string{templateLevel11},
			},
			tplConfig,
			nil,
			fmt.Errorf("more than 10 levels of repository template detected for a-template11->a-template10->a-template9->a-template8->a-template7->a-template6->a-template5->a-template4->a-template3->a-template2->a-template1"),
		},
		"leaf with more than 10 level of template to load": {
			&GhRepoConfig{
				ConfigTemplates: &[]string{TemplateLevel1, templateLevel11, TemplateLevel2},
			},
			tplConfig,
			nil,
			fmt.Errorf("more than 10 levels of repository template detected for a-template11->a-template10->a-template9->a-template8->a-template7->a-template6->a-template5->a-template4->a-template3->a-template2->a-template1"),
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ApplyRepositoryTemplate(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}
