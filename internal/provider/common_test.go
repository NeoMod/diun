package provider

import (
	"testing"

	"github.com/crazy-max/diun/v4/internal/model"
	"github.com/crazy-max/diun/v4/pkg/registry"
	"github.com/crazy-max/diun/v4/pkg/utl"
	"github.com/stretchr/testify/assert"
)

func TestValidateImage(t *testing.T) {
	cases := []struct {
		name          string
		image         string
		metadata      map[string]string
		labels        map[string]string
		watchByDef    bool
		imageDefaults model.Image
		expectedImage model.Image
		expectedErr   interface{}
	}{
		{
			name:       "Test with digest",
			image:      "myimg@sha256:1234567890abcdef",
			watchByDef: true,
			expectedImage: model.Image{
				Name: "myimg@sha256:1234567890abcdef",
			},
			expectedErr: nil,
		},
		{
			name:          "All excluded by default",
			image:         "myimg",
			expectedImage: model.Image{},
			expectedErr:   nil,
		},
		{
			name:       "Include using watch by default",
			image:      "myimg",
			watchByDef: true,
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: nil,
		},
		{
			name:       "Include using diun.enable",
			image:      "myimg",
			watchByDef: false,
			labels: map[string]string{
				"diun.enable": "true",
			},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: nil,
		},
		{
			name:       "Exclude using diun.enable",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.enable": "false",
			},
			expectedImage: model.Image{},
			expectedErr:   nil,
		},
		{
			name:       "Invalid diun.enable",
			image:      "myimg",
			watchByDef: false,
			labels: map[string]string{
				"diun.enable": "chickens",
			},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: &invalidLabelError{},
		},
		// Test diun.regopt
		{
			name:  "Set regopt",
			image: "myimg",
			labels: map[string]string{
				"diun.regopt": "foo",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:   "myimg",
				RegOpt: "foo",
			},
			expectedErr: nil,
		},
		{
			name:  "Set empty recopt",
			image: "myimg",
			labels: map[string]string{
				"diun.regopt": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:   "myimg",
				RegOpt: "",
			},
			expectedErr: nil,
		},
		{
			name:       "Default regopt",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				RegOpt: "foo",
			},
			expectedImage: model.Image{
				Name:   "myimg",
				RegOpt: "foo",
			},
			expectedErr: nil,
		},
		{
			name:       "Override default regopt",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.regopt": "bar",
			},
			imageDefaults: model.Image{
				RegOpt: "foo",
			},
			expectedImage: model.Image{
				Name:   "myimg",
				RegOpt: "bar",
			},
			expectedErr: nil,
		},
		// Test watch_repo
		{
			name:       "Include using global settings",
			image:      "myimg",
			watchByDef: true,
			imageDefaults: model.Image{
				WatchRepo: utl.NewTrue(),
			},
			expectedImage: model.Image{
				Name:      "myimg",
				WatchRepo: utl.NewTrue(),
			},
			expectedErr: nil,
		},
		{
			name:       "Invalid watch_repo",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.watch_repo": "chickens",
			},
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: &invalidLabelError{},
		},
		{
			name:       "Override default image values with labels (true > false)",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.watch_repo": "false",
			},
			imageDefaults: model.Image{
				WatchRepo: utl.NewTrue(),
			},
			expectedImage: model.Image{
				Name:      "myimg",
				WatchRepo: utl.NewFalse(),
			},
			expectedErr: nil,
		},
		{
			name:       "Override default image values with labels (false > true): invalid label error",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.watch_repo": "true",
			},
			imageDefaults: model.Image{
				WatchRepo: utl.NewFalse(),
			},
			expectedImage: model.Image{
				Name:      "myimg",
				WatchRepo: utl.NewTrue(),
			},
			expectedErr: nil,
		},
		// Test diun.notify_on
		{
			name:  "Set valid notify_on",
			image: "myimg",
			labels: map[string]string{
				"diun.notify_on": "new",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:     "myimg",
				NotifyOn: []model.NotifyOn{model.NotifyOnNew},
			},
			expectedErr: nil,
		},
		{
			name:       "Set invalid notify_on",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.notify_on": "chickens",
			},
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:     "myimg",
				NotifyOn: []model.NotifyOn{},
			},
			expectedErr: &invalidLabelError{},
		},
		{
			name:       "Set empty notify_on",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.notify_on": "",
			},
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: nil,
		},
		{
			name:       "Default notify_on",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				NotifyOn: []model.NotifyOn{model.NotifyOnNew},
			},
			expectedImage: model.Image{
				Name:     "myimg",
				NotifyOn: []model.NotifyOn{model.NotifyOnNew},
			},
			expectedErr: nil,
		},
		{
			name:       "Override default notify_on",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.notify_on": "update",
			},
			imageDefaults: model.Image{
				NotifyOn: []model.NotifyOn{model.NotifyOnNew},
			},
			expectedImage: model.Image{
				Name:     "myimg",
				NotifyOn: []model.NotifyOn{model.NotifyOnUpdate},
			},
			expectedErr: nil,
		},
		// Test diun.sort_tags
		{
			name:  "Set valid sort_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.sort_tags": "semver",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:     "myimg",
				SortTags: registry.SortTagSemver,
			},
			expectedErr: nil,
		},
		{
			name:  "Set invalid sort_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.sort_tags": "chickens",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: &invalidLabelError{},
		},
		{
			name:  "Set empty sort_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.sort_tags": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: nil,
		},
		{
			name:       "Default sort_tags",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				SortTags: registry.SortTagSemver,
			},
			expectedImage: model.Image{
				Name:     "myimg",
				SortTags: registry.SortTagSemver,
			},
			expectedErr: nil,
		},
		{
			name:       "Override default sort_tags",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.sort_tags": "reverse",
			},
			imageDefaults: model.Image{
				SortTags: registry.SortTagSemver,
			},
			expectedImage: model.Image{
				Name:     "myimg",
				SortTags: registry.SortTagReverse,
			},
			expectedErr: nil,
		},
		// Test diun.max_tags
		{
			name:  "Set valid max_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.max_tags": "10",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:    "myimg",
				MaxTags: 10,
			},
			expectedErr: nil,
		},
		{
			name:  "Set invalid max_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.max_tags": "chickens",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: &invalidLabelError{},
		},
		{
			name:  "Set empty max_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.max_tags": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: &invalidLabelError{},
		},
		{
			name:       "Default max_tags",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				MaxTags: 10,
			},
			expectedImage: model.Image{
				Name:    "myimg",
				MaxTags: 10,
			},
			expectedErr: nil,
		},
		{
			name:       "Override default max_tags",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.max_tags": "11",
			},
			imageDefaults: model.Image{
				MaxTags: 10,
			},
			expectedImage: model.Image{
				Name:    "myimg",
				MaxTags: 11,
			},
			expectedErr: nil,
		},
		// Test diun.include_tags
		{
			name:  "Set include_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.include_tags": "alpine;ubuntu",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:        "myimg",
				IncludeTags: []string{"alpine", "ubuntu"},
			},
			expectedErr: nil,
		},
		{
			name:  "Set empty include_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.include_tags": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:        "myimg",
				IncludeTags: []string{""},
			},
			expectedErr: nil,
		},
		{
			name:       "Default include_tags",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				IncludeTags: []string{"alpine"},
			},
			expectedImage: model.Image{
				Name:        "myimg",
				IncludeTags: []string{"alpine"},
			},
			expectedErr: nil,
		},
		{
			name:       "Override default include_tags",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.include_tags": "ubuntu",
			},
			imageDefaults: model.Image{
				IncludeTags: []string{"alpine"},
			},
			expectedImage: model.Image{
				Name:        "myimg",
				IncludeTags: []string{"ubuntu"},
			},
			expectedErr: nil,
		},
		// Test diun.exclude_tags
		{
			name:  "Set exclude_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.exclude_tags": "alpine;ubuntu",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:        "myimg",
				ExcludeTags: []string{"alpine", "ubuntu"},
			},
			expectedErr: nil,
		},
		{
			name:  "Set empty exclude_tags",
			image: "myimg",
			labels: map[string]string{
				"diun.exclude_tags": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:        "myimg",
				ExcludeTags: []string{""},
			},
			expectedErr: nil,
		},
		{
			name:       "Default exclude_tags",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				ExcludeTags: []string{"alpine"},
			},
			expectedImage: model.Image{
				Name:        "myimg",
				ExcludeTags: []string{"alpine"},
			},
			expectedErr: nil,
		},
		{
			name:       "Override default exclude_tags",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.exclude_tags": "ubuntu",
			},
			imageDefaults: model.Image{
				ExcludeTags: []string{"alpine"},
			},
			expectedImage: model.Image{
				Name:        "myimg",
				ExcludeTags: []string{"ubuntu"},
			},
			expectedErr: nil,
		},
		// Test diun.hub_tpl
		{
			name:  "Set hub_tpl",
			image: "myimg",
			labels: map[string]string{
				"diun.hub_tpl": "foo",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:   "myimg",
				HubTpl: "foo",
			},
			expectedErr: nil,
		},
		{
			name:  "Set empty hub_tpl",
			image: "myimg",
			labels: map[string]string{
				"diun.hub_tpl": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:   "myimg",
				HubTpl: "",
			},
			expectedErr: nil,
		},
		{
			name:       "Default hub_tpl",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				HubTpl: "foo",
			},
			expectedImage: model.Image{
				Name:   "myimg",
				HubTpl: "foo",
			},
			expectedErr: nil,
		},
		{
			name:       "Override default hub_tpl",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.hub_tpl": "bar",
			},
			imageDefaults: model.Image{
				HubTpl: "foo",
			},
			expectedImage: model.Image{
				Name:   "myimg",
				HubTpl: "bar",
			},
			expectedErr: nil,
		},
		// Test diun.hub_link
		{
			name:  "Set hub_link",
			image: "myimg",
			labels: map[string]string{
				"diun.hub_link": "foo",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:    "myimg",
				HubLink: "foo",
			},
			expectedErr: nil,
		},
		{
			name:  "Set empty hub_link",
			image: "myimg",
			labels: map[string]string{
				"diun.hub_link": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:    "myimg",
				HubLink: "",
			},
			expectedErr: nil,
		},
		{
			name:       "Default hub_link",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				HubLink: "foo",
			},
			expectedImage: model.Image{
				Name:    "myimg",
				HubLink: "foo",
			},
			expectedErr: nil,
		},
		{
			name:       "Override default hub_link",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.hub_link": "bar",
			},
			imageDefaults: model.Image{
				HubLink: "foo",
			},
			expectedImage: model.Image{
				Name:    "myimg",
				HubLink: "bar",
			},
			expectedErr: nil,
		},
		// Test diun.platform
		{
			name:  "Set valid platform",
			image: "myimg",
			labels: map[string]string{
				"diun.platform": "linux/arm/v7",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
				Platform: model.ImagePlatform{
					OS:      "linux",
					Arch:    "arm",
					Variant: "v7",
				},
			},
			expectedErr: nil,
		},
		{
			name:  "Set invalid platform",
			image: "myimg",
			labels: map[string]string{
				"diun.platform": "chickens",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: &invalidLabelError{},
		},
		{
			name:  "Set empty platform",
			image: "myimg",
			labels: map[string]string{
				"diun.platform": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name:     "myimg",
				Platform: model.ImagePlatform{},
			},
			expectedErr: &invalidLabelError{},
		},
		{
			name:       "Default platform",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				Platform: model.ImagePlatform{
					OS:      "linux",
					Arch:    "arm",
					Variant: "v7",
				},
			},
			expectedImage: model.Image{
				Name: "myimg",
				Platform: model.ImagePlatform{
					OS:      "linux",
					Arch:    "arm",
					Variant: "v7",
				},
			},
			expectedErr: nil,
		},
		{
			name:       "Override default platform",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.platform": "linux/arm/v6",
			},
			imageDefaults: model.Image{
				Platform: model.ImagePlatform{
					OS:      "linux",
					Arch:    "arm",
					Variant: "v7",
				},
			},
			expectedImage: model.Image{
				Name: "myimg",
				Platform: model.ImagePlatform{
					OS:      "linux",
					Arch:    "arm",
					Variant: "v6",
				},
			},
			expectedErr: nil,
		},
		// Test diun.metadata
		{
			name:  "Set valid metadata",
			image: "myimg",
			labels: map[string]string{
				"diun.metadata.foo123": "bar",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
				Metadata: map[string]string{
					"foo123": "bar",
				},
			},
			expectedErr: nil,
		},
		{
			name:  "Set invalid metadata",
			image: "myimg",
			labels: map[string]string{
				"diun.metadata.lots of chickens": "bar",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
			expectedErr: &invalidLabelError{},
		},
		{
			name:  "Set empty metadata key",
			image: "myimg",
			labels: map[string]string{
				"diun.metadata.": "bar",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
		},
		{
			name:  "Set empty metadata value",
			image: "myimg",
			labels: map[string]string{
				"diun.metadata.foo123": "",
			},
			watchByDef:    true,
			imageDefaults: model.Image{},
			expectedImage: model.Image{
				Name: "myimg",
			},
		},
		{
			name:       "Default metadata",
			image:      "myimg",
			watchByDef: true,
			labels:     map[string]string{},
			imageDefaults: model.Image{
				Metadata: map[string]string{
					"foo123": "bar",
				},
			},
			expectedImage: model.Image{
				Name: "myimg",
				Metadata: map[string]string{
					"foo123": "bar",
				},
			},
			expectedErr: nil,
		},
		{
			name:       "Merge default metadata",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.metadata.biz123": "baz",
			},
			imageDefaults: model.Image{
				Metadata: map[string]string{
					"foo123": "bar",
				},
			},
			expectedImage: model.Image{
				Name: "myimg",
				Metadata: map[string]string{
					"foo123": "bar",
					"biz123": "baz",
				},
			},
			expectedErr: nil,
		},
		{
			name:       "Override default metadata",
			image:      "myimg",
			watchByDef: true,
			labels: map[string]string{
				"diun.metadata.foo123": "baz",
			},
			imageDefaults: model.Image{
				Metadata: map[string]string{
					"foo123": "bar",
				},
			},
			expectedImage: model.Image{
				Name: "myimg",
				Metadata: map[string]string{
					"foo123": "baz",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			img, err := ValidateImage(
				tt.image,
				tt.metadata,
				tt.labels,
				tt.watchByDef,
				tt.imageDefaults,
			)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedImage, img)
			} else {
				switch err.(type) {
				case *invalidLabelError:
					assert.Error(t, err)
				default:
					assert.Error(t, err)
				}
			}
		})
	}
}
