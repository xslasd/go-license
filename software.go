package client

import (
	"fmt"
	"os"
	"strconv"
)

const (
	OnlineUsers_ItemKey = "online_Users"
	ProgramPath_ItemKey = "program_path"
)

type GetOnlineUsersCallback func() (int64, error)

type onlineUsers struct {
	itemKey  string
	callback GetOnlineUsersCallback
}

func (o onlineUsers) IsLockVal() bool {
	return false
}

func NewOnlineUsers(fn GetOnlineUsersCallback) ActivationHandler {
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
	limit, err := ToInt64E(v)
	if err != nil {
		return ActivationChecksValErr
	}
	if limit <= 0 || limit >= total.(int64) {
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

func ToInt64E(i interface{}) (int64, error) {
	switch s := i.(type) {
	case int64:
		return s, nil
	case int:
		return int64(s), nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to Cast %#v to int64", i)
	case float64:
		return int64(s), nil
	case bool:
		if s {
			return int64(1), nil
		}
		return int64(0), nil
	case nil:
		return int64(0), nil
	default:
		return int64(0), fmt.Errorf("unable to Cast %#v to int64", i)
	}
}
