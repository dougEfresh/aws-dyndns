/*
Copyright © 2019 Douglas Chimento <dchimento@gmail.com>

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
package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"
)

// listRecordsCmd represents the listRecords command
var listRecordsCmd = &cobra.Command{
	Use:   "listRecords",
	Short: "list records of domain",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		zone , err := getZoneId(config.domain)
		if err != nil {
			panic(err)
		}
		query := route53.ListResourceRecordSetsInput{
			HostedZoneId: zone.Id,
		}
		result, err := r53.ListResourceRecordSets(&query)
		if err != nil {
			panic(err)
		}
		for _, value := range result.ResourceRecordSets {
			fmt.Printf("%-40s\t%s\n", *value.Name, *value.ResourceRecords[0].Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(listRecordsCmd)
}
