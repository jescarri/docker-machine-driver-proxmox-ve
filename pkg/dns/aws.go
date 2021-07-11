package dns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/glog"
)

func newSession(awsRegion string, awsRetries int) *session.Session {
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(session.Must(session.NewSession())),
			},
		})
	awsConfig := aws.NewConfig()
	awsConfig.WithCredentials(creds)
	awsConfig.Region = aws.String(awsRegion)
	awsConfig.MaxRetries = aws.Int(awsRetries)

	session, err := session.NewSession(awsConfig)

	if err != nil {
		glog.Errorf("Falied to create aws session: %s", err.Error())
	}

	session.Handlers.Send.PushFront(func(r *request.Request) {
	})

	session.Handlers.Complete.PushFront(func(r *request.Request) {
		if r.Error != nil {
			glog.Infof("Request: %s/%+v, Payload: %s", r.ClientInfo.ServiceName, r.Operation, r.Params)
		}
	})
	return session
}
