package client

import (
	"os"
)

const (
	OnlineUsers_ItemKey = "online_Users"
	ProgramPath_ItemKey = "program_path"
)

type onlineUsers struct {
	itemKey  string
	callback GetValCallback
}

func (o onlineUsers) IsLockVal() bool {
	return false
}

func NewOnlineUsers(fn GetValCallback) ActivationHandler {
	return onlineUsers{itemKey: OnlineUsers_ItemKey, callback: fn}
}

func (o onlineUsers) ItemKey() string {
	return o.itemKey
}

func (o onlineUsers) GetValFn() (any, error) {
	if o.callback != nil {
		return o.callback()
	}
	return -1, nil
}

func (o onlineUsers) CheckFn(data *LicenseInfo, v any) error {
	total, err := o.GetValFn()
	if err != nil {
		return err
	}
	data.NowActivationValues[o.itemKey] = total
	if v == nil || v.(int64) >= total.(int64) {
		return nil
	}
	return OnlineUsersErr
}

type programPath struct {
	itemKey string
}

func NewProgramPath() ActivationHandler {
	return programPath{itemKey: ProgramPath_ItemKey}
}

func (o programPath) IsLockVal() bool {
	return true
}

func (o programPath) ItemKey() string {
	return o.itemKey
}

func (o programPath) GetValFn() (any, error) {
	return os.Getwd()
}

func (o programPath) CheckFn(data *LicenseInfo, v any) error {
	currentPath, err := o.GetValFn()
	if err != nil {
		return err
	}
	if v == nil || v.(string) != currentPath.(string) {
		return LicenseVerifyErr
	}
	return nil
}
