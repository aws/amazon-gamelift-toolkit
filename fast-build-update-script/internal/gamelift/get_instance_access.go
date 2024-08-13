package gamelift

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
)

// InstanceAccessCredentials contains AWS access credentials to get remote SSM access to an instance
type InstanceAccessCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

// GetInstanceAccess will get remote access credentials used to access a GameLift instance using SSM
func (g *GameLiftClient) GetInstanceAccess(ctx context.Context, fleetId string, instanceId string) (*InstanceAccessCredentials, error) {
	getAccessOutput, err := g.gamelift.GetComputeAccess(ctx, &gamelift.GetComputeAccessInput{
		FleetId:     aws.String(fleetId),
		ComputeName: aws.String(instanceId),
	})
	if err != nil {
		return nil, err
	}

	return &InstanceAccessCredentials{
		AccessKeyId:     *getAccessOutput.Credentials.AccessKeyId,
		SecretAccessKey: *getAccessOutput.Credentials.SecretAccessKey,
		SessionToken:    *getAccessOutput.Credentials.SessionToken,
	}, nil
}
