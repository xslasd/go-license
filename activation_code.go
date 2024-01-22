package client

import (
	"encoding/json"
)

type ActivationInfo struct {
	Subject          string         `json:"subject"`
	Description      string         `json:"description,omitempty"`
	InvitationCode   string         `json:"invitation_code"`
	ActivationChecks map[string]any `json:"activation_checks"`
}

type GenerateOption func(*ActivationInfo)

func WithInvitationCode(code string) GenerateOption {
	return func(g *ActivationInfo) {
		g.InvitationCode = code
	}
}

func (c client) GenerateActivationCode(opts ...GenerateOption) ([]byte, error) {
	res := new(ActivationInfo)
	res.ActivationChecks = make(map[string]any)
	for _, o := range opts {
		o(res)
	}
	res.Subject = c.subject
	res.Description = c.description
	for itemKey, item := range c.activationHandlerMap {
		val, err := item.GetValFn()
		if err != nil {
			return nil, err
		}
		switch itemKey {
		default:
			if item.IsLockVal() {
				res.ActivationChecks[itemKey] = val
			}
		}
	}
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return c.activationEncryptFunc(data, c.rsaKey.ActivationEncryptKey)
}
