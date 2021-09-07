// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-licenses/licenses"
	"github.com/spf13/cobra"
)

var (
	checkHelp = "Checks whether licenses for a package are not Forbidden."
	checkCmd  = &cobra.Command{
		Use:   "check <package> [package...]",
		Short: checkHelp,
		Long:  checkHelp + packageHelp,
		Args:  cobra.MinimumNArgs(1),
		RunE:  checkMain,
	}

	excludeRestricted bool
)

func init() {
	checkCmd.Flags().BoolVarP(&excludeRestricted, "exclude-restricted", "", false, "exclude restricted licenses (e.g. GPL)")
	rootCmd.AddCommand(checkCmd)
}

func checkMain(_ *cobra.Command, args []string) error {
	classifier, err := licenses.NewClassifier(confidenceThreshold)
	if err != nil {
		return err
	}

	libs, err := licenses.Libraries(context.Background(), classifier, ignore, args...)
	if err != nil {
		return err
	}

	// indicate that a forbidden license was found
	found := false

	for _, lib := range libs {
		licenseName, licenseType, err := classifier.Identify(lib.LicensePath)
		if err != nil {
			return err
		}

		if excludeRestricted && (licenseType == licenses.Restricted) {
			fmt.Fprintf(os.Stderr, "Restricted license type %s for library %v\n", licenseName, lib)
			found = true
		} else if licenseType == licenses.Forbidden {
			fmt.Fprintf(os.Stderr, "Forbidden license type %s for library %v\n", licenseName, lib)
			found = true
		}
	}

	if found {
		os.Exit(1)
	}

	return nil
}
