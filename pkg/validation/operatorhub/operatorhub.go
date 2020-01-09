package operatorhub

import (
	"fmt"
	"encoding/json"

	"github.com/operator-framework/api/pkg/validation/errors"
	interfaces "github.com/operator-framework/api/pkg/validation/interfaces"
	operatorsv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"

	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

var OperatorHubValidator interfaces.Validator = interfaces.ValidatorFunc(validateOperatorHub)

func validateOperatorHub(objs ...interface{}) (results []errors.ManifestResult) {
	for _, obj := range objs {
		switch v := obj.(type) {
		case *v1alpha1.ClusterServiceVersion:
			results = append(results, validateCSV(v))
		case *registry.ClusterServiceVersion:
			results = append(results, validateCSVRegistry(v))
		}
	}
	return results
}

func validateCSVRegistry(bcsv *registry.ClusterServiceVersion) (result errors.ManifestResult) {
	csv, err := bundleCSVToCSV(bcsv) // internal function
	if err != (errors.Error{}) {
		result.Add(err)
		return result
	}
	return validateCSV(csv)
}

func bundleCSVToCSV(bcsv *registry.ClusterServiceVersion) (*operatorsv1alpha1.ClusterServiceVersion, errors.Error) {
	spec := operatorsv1alpha1.ClusterServiceVersionSpec{}
	if err := json.Unmarshal(bcsv.Spec, &spec); err != nil {
		return nil, errors.ErrInvalidParse(fmt.Sprintf("converting bundle CSV %q", bcsv.GetName()), err)
	}
	return &operatorsv1alpha1.ClusterServiceVersion{
		TypeMeta:   bcsv.TypeMeta,
		ObjectMeta: bcsv.ObjectMeta,
		Spec:       spec,
	}, errors.Error{}
}

// Iterates over the given CSV. Returns a ManifestResult type object.
func validateCSV(csv *v1alpha1.ClusterServiceVersion) errors.ManifestResult {
	result := errors.ManifestResult{Name: csv.GetName()}
	// validate csv for UI.
	result.Add(validateUI(csv)...)
	return result
}

// validateUI validates that the correct fields are populated for proper display in the UI
func validateUI(csv *v1alpha1.ClusterServiceVersion) (errs []errors.Error) {
	annotations := csv.ObjectMeta.GetAnnotations()
	// Return right away if no annotations are found.
	if len(annotations) == 0 {
		errs = append(errs, errors.WarnInvalidCSV("annotations not found", csv.GetName()))
		return errs
	}

	_, descriptionOK := annotations["description"]
	if !descriptionOK {
		errs = append(errs, errors.WarnInvalidCSV("description not found in annotations", csv.GetName()))
		return errs
	}
	// add more checks here
	return errs
}
