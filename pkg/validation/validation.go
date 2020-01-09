package validation

import (
	interfaces "github.com/operator-framework/api/pkg/validation/interfaces"
	operatorhub "github.com/gallettilance/olm-linter/pkg/validation/operatorhub"
)

// OperatorHubValidator implements Validator to validate
// ClusterServiceVersions for the OperatorHub UI.
var OperatorHubValidator = operatorhub.OperatorHubValidator

// AllValidators implements Validator to validate all Operator manifest types.
var AllValidators = interfaces.Validators{
	OperatorHubValidator,
}
