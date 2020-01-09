package operatorhub

import (
	"fmt"
	"os"

	"github.com/gallettilance/olm-linter/pkg/validation"
	"github.com/operator-framework/api/pkg/validation/errors"
	manifests "github.com/operator-framework/api/pkg/manifests"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "operatorhub",
		Short: "Validates all csvs obey operatorhub UI rules",
		Long: `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("command %s requires exactly one argument", cmd.CommandPath())
			}
			_, bundles, _ := manifests.GetManifestsDir(args[0])
			nonEmptyResults := []errors.ManifestResult{}
			objs := []interface{}{}
			for _, obj := range bundles {
				objs = append(objs, obj)
			}
			results := validation.AllValidators.Validate(objs...)
			for _, result := range results {
				if result.HasError() || result.HasWarn() {
					nonEmptyResults = append(nonEmptyResults, result)
				}
			}
			if len(nonEmptyResults) != 0 {
				fmt.Println(nonEmptyResults)
				os.Exit(1)
			}
		},
	}
}