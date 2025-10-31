package cloud

import (
	"fmt"
	"lucky-go/config"
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

// RebootInstance sends a reboot request to the Tencent Cloud platform for the specified destination instance.
// It uses the Tencent Cloud SDK to connect and performs the reboot operation.
func RebootInstance(dest *config.DestinationInstance) error {
	credential := common.NewCredential(os.Getenv("TENCENT_CLOUD_SECRET_ID"), os.Getenv("TENCENT_CLOUD_SECRET_KEY"))

	client, err := lighthouse.NewClient(credential, dest.Region, profile.NewClientProfile())
	if err != nil {
		return err
	}

	response, err := client.RebootInstances(&lighthouse.RebootInstancesRequest{
		InstanceIds: []*string{&dest.InstanceId},
	})

	if err != nil {
		return err
	}

	fmt.Printf("request cloud platform response %v", response.ToJsonString())

	return nil
}
