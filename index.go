package client

import (
	"crypto/sha1"
	"hash"
	"os"
	"path"
	"time"
)

type ActivationHandler interface {
	IsLockVal() bool
	ItemKey() string
	GetValFn() (any, error)
	CheckFn(data *LicenseInfo, v any) error
}

type GetValCallback func() (any, error)

type LicenseCli interface {
	GenerateActivationCode(opts ...GenerateOption) ([]byte, error)
	ActivateLicense(licenseCode []byte) error
	VerifyLicense() bool
	GetLicenseInfo() (*LicenseInfo, error)
}

var pollVerifyTime = time.Hour * 24

type client struct {
	subject         string
	isLockSubject   bool
	description     string
	pollVerifyEvent PollVerifyEvent

	rsaKey RSAKeyConfig

	licenseFileSavePath string
	licenseFileName     string

	activationHandlerMap map[string]ActivationHandler

	activationEncryptFunc ActivationEncryptFunc
	licenseDecryptFunc    LicenseDecryptFunc
	h                     hash.Hash
}
type ActivationEncryptFunc func(plainText []byte, publicKey []byte) ([]byte, error)
type LicenseDecryptFunc func(cipherByte []byte, privateKey []byte) ([]byte, error)

type PollVerifyEvent func(licenseInfo *LicenseInfo, err error)

type Option func(*client)

func WithActivationEncryptFunc(fn ActivationEncryptFunc) Option {
	return func(config *client) {
		config.activationEncryptFunc = fn
	}
}
func WithLicenseDecryptFunc(fn LicenseDecryptFunc) Option {
	return func(config *client) {
		config.licenseDecryptFunc = fn
	}
}
func WithOAEPHash(h hash.Hash) Option {
	return func(config *client) {
		config.h = h
	}
}
func WithPollVerifyEvent(event PollVerifyEvent) Option {
	return func(config *client) {
		config.pollVerifyEvent = event
	}
}

func WithAddActivationHandler(handler ActivationHandler) Option {
	return func(config *client) {
		config.activationHandlerMap[handler.ItemKey()] = handler
	}
}

func WithActivationHandlerMap(handlerMap map[string]ActivationHandler) Option {
	return func(config *client) {
		config.activationHandlerMap = handlerMap
	}
}
func WithIsLockSubject(isLock bool) Option {
	return func(config *client) {
		config.isLockSubject = isLock
	}
}
func WithDescription(description string) Option {
	return func(config *client) {
		config.description = description
	}
}

func WithLicenseFileSavePath(path string) Option {
	return func(config *client) {
		config.licenseFileSavePath = path
	}
}
func WithLicenseFileName(name string) Option {
	return func(config *client) {
		config.licenseFileName = name
	}
}

type RSAKeyConfig struct {
	ActivationEncryptKey []byte
	LicenseDecryptKey    []byte
}

func NewLicenseCli(rsaKey RSAKeyConfig, subject string, opts ...Option) (LicenseCli, error) {
	c := new(client)
	c.subject = subject
	c.rsaKey = rsaKey
	c.licenseFileName = "license.key"
	c.activationHandlerMap = map[string]ActivationHandler{
		SystemOS_ItemKey:    NewSystemOSInfo(),
		CPUInfo_ItemKey:     NewCPUInfo(),
		ProgramPath_ItemKey: NewProgramPath(),
	}
	for _, o := range opts {
		o(c)
	}
	if len(c.activationHandlerMap) == 0 {
		return nil, ActivationHandlerErr
	}
	if c.h == nil {
		c.h = sha1.New()
	}
	if c.activationEncryptFunc == nil {
		c.activationEncryptFunc = c.encrypt
	}
	if c.licenseDecryptFunc == nil {
		c.licenseDecryptFunc = c.decrypt
	}

	if c.licenseFileSavePath == "" {
		c.licenseFileSavePath = c.licenseFileName
	} else {
		err := os.MkdirAll(c.licenseFileSavePath, os.ModeDir)
		if err != nil {
			return nil, err
		}
	}
	c.licenseFileSavePath = path.Join(c.licenseFileSavePath, c.licenseFileName)
	if c.pollVerifyEvent != nil {
		go func() {
			for {
				<-time.After(pollVerifyTime)
				c.pollVerifyEvent(c.GetLicenseInfo())
			}
		}()
	}
	return c, nil
}
