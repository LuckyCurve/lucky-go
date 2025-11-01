package cloud

import (
	"errors"
	"os"
	"testing"

	"lucky-go/config"
)

func TestRebootInstance(t *testing.T) {
	t.Run("SuccessfulReboot", func(t *testing.T) {
		// 保存原始函数
		originalFunc := rebootInstanceFunc
		defer func() {
			rebootInstanceFunc = originalFunc
		}()
		
		// 模拟成功重启
		rebootInstanceFunc = func(dest *config.DestinationInstance) error {
			if dest.Region != "ap-beijing" || dest.InstanceId != "ins-test123" {
				return errors.New("unexpected destination instance values")
			}
			return nil
		}
		
		dest := &config.DestinationInstance{
			Region:     "ap-beijing",
			InstanceId: "ins-test123",
		}
		
		err := RebootInstance(dest)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("RebootError", func(t *testing.T) {
		// 保存原始函数
		originalFunc := rebootInstanceFunc
		defer func() {
			rebootInstanceFunc = originalFunc
		}()
		
		// 模拟重启失败
		expectedErr := errors.New("reboot failed")
		rebootInstanceFunc = func(dest *config.DestinationInstance) error {
			return expectedErr
		}
		
		dest := &config.DestinationInstance{
			Region:     "ap-shanghai",
			InstanceId: "ins-test456",
		}
		
		err := RebootInstance(dest)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected error '%v', got '%v'", expectedErr, err)
		}
	})

	t.Run("WithEnvironmentVariables", func(t *testing.T) {
		// 保存原始函数
		originalFunc := rebootInstanceFunc
		defer func() {
			rebootInstanceFunc = originalFunc
		}()
		
		// 设置环境变量
		originalSecretID := os.Getenv("TENCENT_CLOUD_SECRET_ID")
		originalSecretKey := os.Getenv("TENCENT_CLOUD_SECRET_KEY")
		os.Setenv("TENCENT_CLOUD_SECRET_ID", "test_secret_id")
		os.Setenv("TENCENT_CLOUD_SECRET_KEY", "test_secret_key")
		defer func() {
			os.Setenv("TENCENT_CLOUD_SECRET_ID", originalSecretID)
			os.Setenv("TENCENT_CLOUD_SECRET_KEY", originalSecretKey)
		}()
		
		// 模拟使用环境变量的重启函数
		rebootInstanceFunc = func(dest *config.DestinationInstance) error {
			secretID := os.Getenv("TENCENT_CLOUD_SECRET_ID")
			secretKey := os.Getenv("TENCENT_CLOUD_SECRET_KEY")
			
			if secretID != "test_secret_id" || secretKey != "test_secret_key" {
				return errors.New("environment variables not set correctly")
			}
			
			if dest.Region != "ap-guangzhou" || dest.InstanceId != "ins-test789" {
				return errors.New("unexpected destination instance values")
			}
			return nil
		}
		
		dest := &config.DestinationInstance{
			Region:     "ap-guangzhou",
			InstanceId: "ins-test789",
		}
		
		err := RebootInstance(dest)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}