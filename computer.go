package client

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"strings"
)

const (
	SystemOS_ItemKey = "system_os"
	CPUInfo_ItemKey  = "cpu_info"
)

type systemOSInfo struct {
	itemKey string
}

func (s systemOSInfo) IsLockVal() bool {
	return true
}

func NewSystemOSInfo() ActivationHandler {
	return systemOSInfo{itemKey: SystemOS_ItemKey}
}

func (s systemOSInfo) ItemKey() string {
	return s.itemKey
}

func (s systemOSInfo) GetValFn() (a any, err error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s(%s):%s", info.Platform, info.Hostname, info.HostID), nil
}

func (s systemOSInfo) CheckFn(data *LicenseInfo, v any) error {
	hostStr, err := s.GetValFn()
	if err != nil {
		return err
	}
	if v == nil || v.(string) != hostStr.(string) {
		return LicenseVerifyErr
	}
	return nil
}

type CPUInfo struct {
	itemKey string
}

func (c CPUInfo) IsLockVal() bool {
	return true
}

func NewCPUInfo() ActivationHandler {
	return CPUInfo{itemKey: CPUInfo_ItemKey}
}
func (c CPUInfo) ItemKey() string {
	return c.itemKey
}

func (c CPUInfo) GetValFn() (a any, err error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	if len(info) == 0 {
		return "not cpu info", err
	}
	cpu0 := info[0]
	return fmt.Sprintf("%s:%s", strings.TrimSpace(cpu0.ModelName), cpu0.PhysicalID), nil
}

func (c CPUInfo) CheckFn(data *LicenseInfo, v any) error {
	cpuStr, err := c.GetValFn()
	if err != nil {
		return err
	}
	if v == nil || v.(string) != cpuStr.(string) {
		return LicenseVerifyErr
	}
	return nil
}
