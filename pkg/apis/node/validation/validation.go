/*
Copyright 2019 The Kubernetes Authors.

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

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/apis/core"
	corevalidation "k8s.io/kubernetes/pkg/apis/core/validation"
	"k8s.io/kubernetes/pkg/apis/node"
	"k8s.io/kubernetes/pkg/features"
)

// ValidateRuntimeClass validates the RuntimeClass
func ValidateRuntimeClass(rc *node.RuntimeClass) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&rc.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	for _, msg := range apivalidation.NameIsDNSLabel(rc.Handler, false) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("handler"), rc.Handler, msg))
	}

	if rc.Overhead != nil && utilfeature.DefaultFeatureGate.Enabled(features.PodOverhead) {
		allErrs = append(allErrs, validateOverhead(rc.Overhead, field.NewPath("overhead"))...)
	}

	return allErrs
}

// ValidateRuntimeClassUpdate validates an update to the object
func ValidateRuntimeClassUpdate(new, old *node.RuntimeClass) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&new.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))

	allErrs = append(allErrs, apivalidation.ValidateImmutableField(new.Handler, old.Handler, field.NewPath("handler"))...)

	return allErrs
}

func validateOverhead(overhead *node.Overhead, fldPath *field.Path) field.ErrorList {
	// reuse the ResourceRequirements validation logic
	return corevalidation.ValidateResourceRequirements(&core.ResourceRequirements{Limits: overhead.PodFixed}, fldPath)
}
