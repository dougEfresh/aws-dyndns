// Copyright © 2019.  Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var r53 *route53.Route53

func setupAws() error {
	sess, err := session.NewSessionWithOptions(session.Options{})
	if err != nil {
		return err
	}
	r53 = route53.New(sess)
	return nil
}

func getZoneId(domain string) (*route53.HostedZone, error) {
	query := &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(config.domain),
	}
	zone, err := r53.ListHostedZonesByName(query)
	if err != nil {
		return nil, err
	}
	logger.Debug(fmt.Sprintf("Got back %p %s %s", zone.HostedZoneId, *zone.HostedZones[0].Id, *zone.HostedZones[0].Name))
	return zone.HostedZones[0], nil
}
