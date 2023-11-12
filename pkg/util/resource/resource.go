/*
Copyright 2021 Vesoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resource

import (
	kruisev1beta1 "github.com/openkruise/kruise-api/apps/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vesoft-inc/nebula-operator/apis/apps/v1alpha1"
	"github.com/vesoft-inc/nebula-operator/pkg/util/discovery"
)

type GVRFunc func() schema.GroupVersionResource

var (
	NebulaClusterKind       = v1alpha1.SchemeGroupVersion.WithKind("NebulaCluster")
	StatefulSetKind         = appsv1.SchemeGroupVersion.WithKind("StatefulSet")
	AdvancedStatefulSetKind = kruisev1beta1.SchemeGroupVersion.WithKind("StatefulSet")

	GroupVersionResources = map[string]GVRFunc{
		StatefulSetKind.String(): GetStatefulSetGVR,
	}
)

func GetStatefulSetGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "statefulsets",
	}
}

func GetAdvancedStatefulSetGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "apps.kruise.io",
		Version:  "v1beta1",
		Resource: "statefulsets",
	}
}

func GetUniteDeploymentGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "apps.kruise.io",
		Version:  "v1alpha1",
		Resource: "uniteddeployments",
	}
}

func GetGVKFromDefinition(dm discovery.Interface, ref v1alpha1.WorkloadReference) (schema.GroupVersionKind, error) {
	// if given definitionRef is empty return a default GVK
	if ref.Name == "" {
		return StatefulSetKind, nil
	}
	var gvk schema.GroupVersionKind
	groupResource := schema.ParseGroupResource(ref.Name)
	gvr := schema.GroupVersionResource{Group: groupResource.Group, Resource: groupResource.Resource, Version: ref.Version}
	kinds, err := dm.KindsFor(gvr)
	if err != nil {
		return gvk, err
	}
	if len(kinds) < 1 {
		return gvk, &meta.NoResourceMatchError{
			PartialResource: gvr,
		}
	}
	return kinds[0], nil
}
