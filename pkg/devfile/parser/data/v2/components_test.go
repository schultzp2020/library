package v2

import (
	"reflect"
	"testing"

	v1 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/api/v2/pkg/attributes"
	"github.com/devfile/library/pkg/devfile/parser/data/v2/common"
	"github.com/devfile/library/pkg/testingutil"
	"github.com/stretchr/testify/assert"
)

func TestDevfile200_AddComponent(t *testing.T) {

	tests := []struct {
		name              string
		currentComponents []v1.Component
		newComponents     []v1.Component
		wantErr           bool
	}{
		{
			name: "successfully add the component",
			currentComponents: []v1.Component{
				{
					Name: "component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "component2",
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			newComponents: []v1.Component{
				{
					Name: "component3",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error out on duplicate component",
			currentComponents: []v1.Component{
				{
					Name: "component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "component2",
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			newComponents: []v1.Component{
				{
					Name: "component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.currentComponents,
						},
					},
				},
			}

			err := d.AddComponents(tt.newComponents)
			// Unexpected error
			if (err != nil) != tt.wantErr {
				t.Errorf("TestDevfile200_AddComponents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDevfile200_UpdateComponent(t *testing.T) {

	tests := []struct {
		name              string
		currentComponents []v1.Component
		newComponent      v1.Component
		wantErr           bool
	}{
		{
			name: "successfully update the component",
			currentComponents: []v1.Component{
				{
					Name: "Component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image1",
							},
						},
					},
				},
				{
					Name: "component2",
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			newComponent: v1.Component{
				Name: "Component1",
				ComponentUnion: v1.ComponentUnion{
					Container: &v1.ContainerComponent{
						Container: v1.Container{
							Image: "image2",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail to update the component if not exist",
			currentComponents: []v1.Component{
				{
					Name: "Component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image1",
							},
						},
					},
				},
			},
			newComponent: v1.Component{
				Name: "Component2",
				ComponentUnion: v1.ComponentUnion{
					Container: &v1.ContainerComponent{
						Container: v1.Container{
							Image: "image2",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.currentComponents,
						},
					},
				},
			}

			err := d.UpdateComponent(tt.newComponent)
			// Unexpected error
			if (err != nil) != tt.wantErr {
				t.Errorf("TestDevfile200_UpdateComponent() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				components, err := d.GetComponents(common.DevfileOptions{})
				if err != nil {
					t.Errorf("TestDevfile200_UpdateComponent() unxpected error %v", err)
					return
				}

				matched := false
				for _, component := range components {
					if reflect.DeepEqual(component, tt.newComponent) {
						matched = true
						break
					}
				}

				if !matched {
					t.Error("TestDevfile200_UpdateComponent() error updating the component")
				}
			}
		})
	}
}

func TestGetDevfileComponents(t *testing.T) {

	tests := []struct {
		name           string
		component      []v1.Component
		wantComponents []string
		filterOptions  common.DevfileOptions
		wantErr        bool
	}{
		{
			name:      "Invalid devfile",
			component: []v1.Component{},
		},
		{
			name: "Get all the components",
			component: []v1.Component{
				{
					Name: "comp1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"secondString": "secondStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "comp2",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"fourthString": "fourthStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			wantComponents: []string{"comp1", "comp2"},
		},
		{
			name: "Get component with the specified filter",
			component: []v1.Component{
				{
					Name: "comp1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"secondString": "secondStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "comp2",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "comp3",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"fourthString": "fourthStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
				{
					Name: "comp4",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"fourthString": "fourthStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstString": "firstStringValue",
				},
				CommandOptions: common.CommandOptions{
					CommandGroupKind: v1.BuildCommandGroupKind,
					CommandType:      v1.CompositeCommandType,
				},
				ComponentOptions: common.ComponentOptions{
					ComponentType: v1.VolumeComponentType,
				},
			},
			wantComponents: []string{"comp3"},
		},
		{
			name: "Wrong filter for component",
			component: []v1.Component{
				{
					Name: "comp1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"secondString": "secondStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "comp2",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstStringIsWrong": "firstStringValue",
				},
				ComponentOptions: common.ComponentOptions{
					ComponentType: v1.ContainerComponentType,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid component type",
			component: []v1.Component{
				{
					Name: "comp1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstString": "firstStringValue",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.component,
						},
					},
				},
			}

			components, err := d.GetComponents(tt.filterOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestGetDevfileComponents() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				// confirm the length of actual vs expected
				if len(components) != len(tt.wantComponents) {
					t.Errorf("TestGetDevfileComponents() error - length of expected components is not the same as the length of actual components")
					return
				}

				// compare the component slices for content
				for _, wantComponent := range tt.wantComponents {
					matched := false
					for _, component := range components {
						if wantComponent == component.Name {
							matched = true
						}
					}

					if !matched {
						t.Errorf("TestGetDevfileComponents() error - component %s not found in the devfile", wantComponent)
					}
				}
			}
		})
	}

}

func TestGetDevfileContainerComponents(t *testing.T) {

	tests := []struct {
		name                 string
		component            []v1.Component
		expectedMatchesCount int
		filterOptions        common.DevfileOptions
		wantErr              bool
	}{
		{
			name:                 "Invalid devfile",
			component:            []v1.Component{},
			expectedMatchesCount: 0,
		},
		{
			name: "Valid devfile with wrong component type (Openshift)",
			component: []v1.Component{
				{
					ComponentUnion: v1.ComponentUnion{
						Openshift: &v1.OpenshiftComponent{},
					},
				},
			},
			expectedMatchesCount: 0,
		},
		{
			name: "Valid devfile with correct component type (Container)",
			component: []v1.Component{
				testingutil.GetFakeContainerComponent("comp1"),
				testingutil.GetFakeContainerComponent("comp2"),
			},
			expectedMatchesCount: 2,
			filterOptions:        common.DevfileOptions{},
		},
		{
			name: "Get Container component with the specified filter",
			component: []v1.Component{
				{
					Name: "comp1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"secondString": "secondStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "comp2",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstString":  "firstStringValue",
					"secondString": "secondStringValue",
				},
			},
			expectedMatchesCount: 1,
		},
		{
			name: "Get Container component with the wrong specified filter",
			component: []v1.Component{
				{
					Name: "comp1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"secondString": "secondStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "comp2",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstStringIsWrong": "firstStringValue",
				},
			},
			expectedMatchesCount: 0,
			wantErr:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.component,
						},
					},
				},
			}

			devfileComponents, err := d.GetDevfileContainerComponents(tt.filterOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestGetDevfileContainerComponents() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && len(devfileComponents) != tt.expectedMatchesCount {
				t.Errorf("TestGetDevfileContainerComponents error: wrong number of components matched: expected %v, actual %v", tt.expectedMatchesCount, len(devfileComponents))
			}
		})
	}

}

func TestGetDevfileVolumeComponents(t *testing.T) {

	tests := []struct {
		name                 string
		component            []v1.Component
		expectedMatchesCount int
		filterOptions        common.DevfileOptions
		wantErr              bool
	}{
		{
			name:                 "Invalid devfile",
			component:            []v1.Component{},
			expectedMatchesCount: 0,
		},
		{
			name: "Valid devfile with wrong component type (Kubernetes)",
			component: []v1.Component{
				{
					ComponentUnion: v1.ComponentUnion{
						Kubernetes: &v1.KubernetesComponent{},
					},
				},
			},
			expectedMatchesCount: 0,
		},
		{
			name: "Valid devfile with correct component type (Volume)",
			component: []v1.Component{
				testingutil.GetFakeContainerComponent("comp1"),
				testingutil.GetFakeVolumeComponent("myvol", "4Gi"),
			},
			expectedMatchesCount: 1,
		},
		{
			name: "Get Container component with the specified filter",
			component: []v1.Component{
				{
					Name: "comp1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"secondString": "secondStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
				{
					Name: "comp2",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString": "firstStringValue",
						"thirdString": "thirdStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstString": "firstStringValue",
				},
			},
			expectedMatchesCount: 2,
		},
		{
			name: "Get Container component with the wrong specified filter",
			component: []v1.Component{
				{
					Name: "comp1",
					Attributes: attributes.Attributes{}.FromStringMap(map[string]string{
						"firstString":  "firstStringValue",
						"secondString": "secondStringValue",
					}),
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			filterOptions: common.DevfileOptions{
				Filter: map[string]interface{}{
					"firstStringIsWrong": "firstStringValue",
				},
			},
			expectedMatchesCount: 0,
			wantErr:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.component,
						},
					},
				},
			}
			devfileComponents, err := d.GetDevfileVolumeComponents(tt.filterOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestGetDevfileVolumeComponents() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && len(devfileComponents) != tt.expectedMatchesCount {
				t.Errorf("TestGetDevfileVolumeComponents error: wrong number of components matched: expected %v, actual %v", tt.expectedMatchesCount, len(devfileComponents))
			}
		})
	}

}

func TestDeleteComponents(t *testing.T) {

	d := &DevfileV2{
		v1.Devfile{
			DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
				DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
					Components: []v1.Component{
						{
							Name: "comp2",
							ComponentUnion: v1.ComponentUnion{
								Container: &v1.ContainerComponent{
									Container: v1.Container{
										VolumeMounts: []v1.VolumeMount{
											testingutil.GetFakeVolumeMount("comp2", "/path"),
											testingutil.GetFakeVolumeMount("comp2", "/path2"),
											testingutil.GetFakeVolumeMount("comp3", "/path"),
										},
									},
								},
							},
						},
						{
							Name: "comp2",
							ComponentUnion: v1.ComponentUnion{
								Volume: &v1.VolumeComponent{},
							},
						},
						{
							Name: "comp3",
							ComponentUnion: v1.ComponentUnion{
								Volume: &v1.VolumeComponent{},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name              string
		componentToDelete string
		wantComponents    []v1.Component
		wantErr           bool
	}{
		{
			name:              "Successfully delete a Component",
			componentToDelete: "comp3",
			wantComponents: []v1.Component{
				{
					Name: "comp2",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								VolumeMounts: []v1.VolumeMount{
									testingutil.GetFakeVolumeMount("comp2", "/path"),
									testingutil.GetFakeVolumeMount("comp2", "/path2"),
									testingutil.GetFakeVolumeMount("comp3", "/path"),
								},
							},
						},
					},
				},
				{
					Name: "comp2",
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			wantErr: false,
		},
		{
			name:              "Missing Component",
			componentToDelete: "comp12",
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := d.DeleteComponent(tt.componentToDelete)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteComponent() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				assert.Equal(t, tt.wantComponents, d.Components, "The two values should be the same.")
			}
		})
	}

}
