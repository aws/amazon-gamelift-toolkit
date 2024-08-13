package gamelift

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	"github.com/aws/aws-sdk-go-v2/service/gamelift/types"
)

// OpenPortForFleet will open the provided port for the range of IPs provided
func (g *GameLiftClient) OpenPortForFleet(ctx context.Context, fleetId string, port int32, ipRange string) error {
	_, err := g.gamelift.UpdateFleetPortSettings(ctx, &gamelift.UpdateFleetPortSettingsInput{
		FleetId: aws.String(fleetId),
		InboundPermissionAuthorizations: []types.IpPermission{
			types.IpPermission{
				FromPort: aws.Int32(port),
				IpRange:  aws.String(ipRange),
				Protocol: types.IpProtocolTcp,
				ToPort:   aws.Int32(port),
			},
		},
	})

	// If we have already opened this port on this fleet, there is no reason to return an error
	if err != nil {
		ire := new(types.InvalidRequestException)
		if errors.As(err, &ire) {
			if strings.Contains(ire.ErrorMessage(), "InvalidPermission.Duplicate") {
				err = nil
			}
		}
	}

	return err
}
