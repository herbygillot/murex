//go:build plan9
// +build plan9

package defaults

/*
   WARNING:
   --------

   This Go source file has been automatically generated from
   profile_plan9.mx using docgen.

   Please do not manually edit this file because it will be automatically
   overwritten by the build pipeline. Instead please edit the aforementioned
   profile_plan9.mx file located in the same directory.
*/

func init() {
	murexProfile = append(murexProfile, "autocomplete: set bg {\n    [{\n        \"DynamicDesc\": ({ fid-list --stopped }),\n        \"ListView\": true\n    }]\n}\n\nautocomplete: set fg {\n    [{\n        \"DynamicDesc\": ({ fid-list --stopped }),\n        \"ListView\": true\n    }]\n}")
}
