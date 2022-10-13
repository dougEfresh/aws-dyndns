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
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "upsert",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var recordValue = args[0]
		ip := net.ParseIP(recordValue)
		addressType := route53.RRTypeAaaa
		if ip.To4() != nil {
			addressType = route53.RRTypeA
		}

		zone, err := getZoneId(config.domain)
		if err != nil {
			panic(err)
		}
		input := &route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &route53.ChangeBatch{
				Changes: []*route53.Change{
					{
						Action: aws.String(route53.ChangeActionUpsert),
						ResourceRecordSet: &route53.ResourceRecordSet{
							Name: aws.String(fmt.Sprintf("%s.%s", config.record, config.domain)),
							ResourceRecords: []*route53.ResourceRecord{
								{
									Value: aws.String(recordValue),
								},
							},
							TTL:  aws.Int64(600),
							Type: aws.String(addressType),
						},
					},
				},
				Comment: aws.String(fmt.Sprintf("Update %s ", time.Now())),
			},
			HostedZoneId: zone.Id,
		}
		_, err = r53.ChangeResourceRecordSets(input)
		if err != nil {
			panic(err)
		}
		logger.Info("Updated aws zone with new address", zap.String("address", config.record))
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().StringVarP(&config.record, "record", "r", "", "record")
}
