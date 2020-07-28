//   Copyright 2019 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package riEc2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/trackit/trackit/es/indexes/common"
	"github.com/trackit/trackit/es/indexes/riEc2Reports"
)

func getRecurringCharges(reservation *ec2.ReservedInstances) []riEc2Reports.RecurringCharges {
	charges := make([]riEc2Reports.RecurringCharges, len(reservation.RecurringCharges))
	for i, key := range reservation.RecurringCharges {
		charges[i] = riEc2Reports.RecurringCharges{
			Amount:    aws.Float64Value(key.Amount),
			Frequency: aws.StringValue(key.Frequency),
		}
	}
	return charges
}

// getReservationTag formats []*ec2.Tag to map[string]string
func getReservationTag(tags []*ec2.Tag) []common.Tag {
	res := make([]common.Tag, 0)
	for _, tag := range tags {
		res = append(res, common.Tag{
			Key:   aws.StringValue(tag.Key),
			Value: aws.StringValue(tag.Value),
		})
	}
	return res
}
