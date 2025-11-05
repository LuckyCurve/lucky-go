package cloud

import (
	"fmt"
	"lucky-go/config"
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

// 定义函数变量，用于在测试中模拟
var rebootInstanceFunc = defaultRebootInstance

// RebootInstance 向腾讯云平台发送重启请求以重启指定的目标实例。
// 它使用腾讯云SDK进行连接并执行重启操作。
func RebootInstance(dest *config.DestinationInstance) error {
	return rebootInstanceFunc(dest)
}

// defaultRebootInstance 是 RebootInstance 的默认实现
func defaultRebootInstance(dest *config.DestinationInstance) error {
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

	fmt.Printf("请求云平台响应 %v", response.ToJsonString())

	return nil
}
