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
	"context"
	"fmt"
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		sync()
	},
}

func sync() {
	resolverIps, err := net.LookupIP(config.resolver)
	if err != nil {
		panic(err)
	}
	openDns := buildResolver(resolverIps[0].String())
	nameservers, err := openDns.LookupNS(context.Background(), config.domain)
	if err != nil {
		panic(err)
	}
	awsNs := nameservers[0]
	awsNsIpaddrs, err := openDns.LookupIPAddr(context.Background(), awsNs.Host)
	if err != nil {
		panic(err)
	}
	awsDns := buildResolver(awsNsIpaddrs[0].IP.String())
	names, err := awsDns.LookupIPAddr(context.Background(), fmt.Sprintf("%s.%s", config.record, config.domain))
	if err != nil {
		panic(err)
	}
	awsIp := names[0]
	names, err = openDns.LookupIPAddr(context.Background(), "myip.opendns.com")
	if err != nil {
		panic(err)
	}
	myIp := names[0]
	if myIp.IP.Equal(awsIp.IP) {
		logger.Info(fmt.Sprintf("Same IP as aws %s", myIp.IP))
		return
	}
	logger.Info(fmt.Sprintf("MyIp(%s)!=AwsIp(%s)", myIp.IP, awsIp.IP))
	query := &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(config.domain),
	}
	zone, err := r53.ListHostedZonesByName(query)
	if err != nil {
		panic(err)
	}
	logger.Info(fmt.Sprintf("HostedZoneId:%s", zone.HostedZones[0].Id))
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(fmt.Sprintf("%s.%s", config.record, config.domain)),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(myIp.IP.String()),
							},
						},
						TTL:  aws.Int64(600),
						Type: aws.String("A"),
					},
				},
			},
			Comment: aws.String(fmt.Sprintf("previous %s ", awsIp.IP)),
		},
		HostedZoneId: zone.HostedZones[0].Id,
	}

	_, err = r53.ChangeResourceRecordSets(input)
	if err != nil {
		panic(err)
	}
	logger.Info("Updated aws zone with new address")
}


func buildResolver(resolver string) net.Resolver {
	logger.Debug(fmt.Sprintf("Using %s for resolver", resolver))
	return net.Resolver{
		PreferGo:     true,
		StrictErrors: false,
		Dial: builDialer(resolver),
	}
}

func builDialer(resolver string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (conn net.Conn, e error) {
		d := net.Dialer{
			Timeout: time.Millisecond * time.Duration(1000),
		}
		r := fmt.Sprintf("%s:53", resolver)
		logger.Debug(fmt.Sprintf("Calling %s", r))
		return d.DialContext(ctx, "udp", r)
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.PersistentFlags().StringVar(&config.record, "record", "", "record (type A) of domain")
	syncCmd.PersistentFlags().StringVar(&config.resolver, "resolver", "resolver1.opendns.com", "use resolver to lookup address")
}
