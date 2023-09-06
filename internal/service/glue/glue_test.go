// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package glue_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccGlue_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"DataCatalogEncryptionSettings": {
			"basic":      testAccDataCatalogEncryptionSettings_basic,
			"dataSource": testAccDataCatalogEncryptionSettingsDataSource_basic,
		},
		"ResourcePolicy": {
			"basic":      testAccResourcePolicy_basic,
			"update":     testAccResourcePolicy_update,
			"hybrid":     testAccResourcePolicy_hybrid,
			"disappears": testAccResourcePolicy_disappears,
			"equivalent": testAccResourcePolicy_ignoreEquivalent,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
