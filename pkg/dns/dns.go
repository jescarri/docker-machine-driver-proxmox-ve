package dns

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/golang/glog"
)

func CreateRecord(ipAddress string, hostname string, hostedZoneId string, ttl int64, weight int64, awsRegion string) error {
	return r53Call("UPSERT", ipAddress, hostname, hostedZoneId, ttl, weight, awsRegion)
}

func DeleteRecord(ipAddress string, hostname string, hostedZoneId string, ttl int64, weight int64, awsRegion string) error {
	return r53Call("DELETE", ipAddress, hostname, hostedZoneId, ttl, weight, awsRegion)
}

func r53Call(action string, ipAddress string, hostname string, hostedZoneId string, ttl int64, weight int64, awsRegion string) error {
	sess := newSession(awsRegion, 5)
	svc := route53.New(sess)

	input := &route53.GetHostedZoneInput{
		Id: aws.String(hostedZoneId),
	}
	result, err := svc.GetHostedZone(input)
	if err != nil {
		return err
	}
	domain := *result.HostedZone.Name
	fqdn := fmt.Sprintf("%s.%s", hostname, domain)
	glog.Infof("Adding DNS Record for Host: %s with Ip: %s", fqdn, ipAddress)
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{ // Required
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String(action), // Required
					ResourceRecordSet: &route53.ResourceRecordSet{ // Required
						Name: aws.String(fqdn), // Required
						Type: aws.String("A"),  // Required
						ResourceRecords: []*route53.ResourceRecord{
							{ // Required
								Value: aws.String(ipAddress), // Required
							},
						},
						TTL:           aws.Int64(ttl),
						Weight:        aws.Int64(weight),
						SetIdentifier: aws.String("managed by proxmoxve-docker-machine"),
					},
				},
			},
			Comment: aws.String("managed by proxmoxve-docker-machine"),
		},
		HostedZoneId: aws.String(hostedZoneId), // Required
	}
	resp, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		glog.Infof("%s", resp)
		return err
	}
	return nil

}
