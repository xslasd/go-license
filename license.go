package client

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type LicenseInfo struct {
	LicenseKey     string        `json:"license_key" form:"license_key"`
	Subject        string        `json:"subject"`
	Description    string        `json:"description,omitempty"`
	IssuedTime     int64         `json:"issued_time"`
	ExpiryTime     int64         `json:"expiry_time"`
	InvitationCode string        `json:"invitation_code,omitempty"`
	PollVerifyTime time.Duration `json:"poll_verify_time"`

	ActivationChecks    map[string]any `json:"activation_checks"`
	NowActivationValues map[string]any `json:"now_activation_values"`
}

func (c client) VerifyLicense() bool {
	_, err := c.GetLicenseInfo()
	if err != nil {
		return false
	}
	return true
}

func (c client) ActivateLicense(licenseCode []byte) (*LicenseInfo, error) {
	data, err := c.getServerLicenseInfo(licenseCode)
	if err != nil {
		return nil, err
	}
	err = c.verify(data)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(c.licenseFileSavePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Write(licenseCode)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c client) GetLicenseInfo() (*LicenseInfo, error) {
	licenseCode, err := os.ReadFile(c.licenseFileSavePath)
	if err != nil {
		return nil, err
	}
	data, err := c.getServerLicenseInfo(licenseCode)
	if err != nil {
		return nil, err
	}
	err = c.verify(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (c client) verify(data *LicenseInfo) error {
	if c.isLockSubject && c.subject != data.Subject {
		return LicenseVerifyErr
	}
	if data.ExpiryTime > -1 {
		t := time.Now()
		if c.lastRecordedTimeFunc != nil {
			lt := c.lastRecordedTimeFunc()
			if lt.After(t) {
				return LicenseExpirationErr
			}
		}
		e := time.UnixMilli(data.ExpiryTime)
		if t.After(e) {
			return LicenseExpirationErr
		}
		data.NowActivationValues["end_time"] = e.Sub(t)
	}
	for itemKey, item := range c.activationHandlerMap {
		v := data.ActivationChecks[itemKey]
		err := item.CheckFn(data, v)
		if err != nil {
			return err
		}
	}
	fmt.Println("License verification successful!")
	return nil
}
func (c client) getServerLicenseInfo(licenseCode []byte) (*LicenseInfo, error) {
	res, err := c.licenseDecryptFunc(licenseCode, c.rsaKey.LicenseDecryptKey)
	if err != nil {
		return nil, err
	}
	data := new(LicenseInfo)
	err = json.Unmarshal(res, data)
	data.NowActivationValues = make(map[string]any)
	return data, err
}
