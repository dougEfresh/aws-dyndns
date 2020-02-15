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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a  record",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		zone, err := getZoneId(config.domain)
		if err != nil {
			panic(err)
		}
		query := route53.ListResourceRecordSetsInput{
			HostedZoneId:    zone.Id,
			StartRecordName: aws.String(fmt.Sprintf("%s.%s", config.record, config.domain)),
		}
		records, err := r53.ListResourceRecordSets(&query)
		if err != nil {
			panic(err)
		}
		if len(records.ResourceRecordSets) == 0 {
			panic(fmt.Errorf("no records for %s.%s", config.record, config.domain))
		}
		resourceRecords := records.ResourceRecordSets
		logger.Debug(fmt.Sprintf("Found %d records", len(resourceRecords)))
		var changeSet = make([]*route53.ResourceRecordSet, len(resourceRecords))

		for _, r := range resourceRecords {
			changeSet = append(changeSet, &route53.ResourceRecordSet{
				Name: r.Name,
				TTL:  r.TTL,
			})
		}

		input := &route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &route53.ChangeBatch{
				Changes: []*route53.Change{
					{
						Action:            aws.String(route53.ChangeActionDelete),
						ResourceRecordSet: resourceRecords[0],
					},
				},
			},
			HostedZoneId: zone.Id,
		}
		_, err = r53.ChangeResourceRecordSets(input)
		if err != nil {
			panic(err)
		}
		logger.Info("Deleted", zap.String("address", config.record))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringVarP(&config.record, "record", "r", "", "record")
}
