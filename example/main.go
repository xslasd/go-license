package main

import (
	"flag"
	"fmt"
	client "go-license"
	"os"
)

func main() {
	var generateActivationCode bool
	flag.BoolVar(&generateActivationCode, "GenerateActivationCode", false, "Generate activation file")
	var invitationCode string
	flag.StringVar(&invitationCode, "InvitationCode", "", "Set activate invitation code")

	var verifyLicense bool
	flag.BoolVar(&verifyLicense, "VerifyLicense", false, "Verify License")

	var activationEncryptKeyDir string
	var licenseDecryptKeyDir string
	flag.StringVar(&activationEncryptKeyDir, "ActivationEncryptKeyDir", "./.client_key/activation_encrypt.pem", "Set encrypt activation public key file directory")
	flag.StringVar(&licenseDecryptKeyDir, "LicenseDecryptKeyDir", "./.client_key/license_decrypt.pem", "Set decryption license private key file directory")
	flag.Parse()

	activationEncryptKey, err := os.ReadFile(activationEncryptKeyDir)
	if err != nil {
		panic(err)
	}

	licenseDecryptKey, err := os.ReadFile(licenseDecryptKeyDir)
	if err != nil {
		panic(err)
	}

	cli, err := client.NewLicenseCli(client.RSAKeyConfig{
		ActivationEncryptKey: activationEncryptKey,
		LicenseDecryptKey:    licenseDecryptKey,
	}, "Demo License")
	if err != nil {
		panic(err)
	}

	if generateActivationCode {
		opts := make([]client.GenerateOption, 0)
		if invitationCode != "" {
			opts = append(opts, client.WithInvitationCode(invitationCode))
		}
		code, err := cli.GenerateActivationCode(opts...)
		if err != nil {
			panic(err)
		}

		f, err := os.Create("activation_code.key")
		if err != nil {
			panic(err)
		}
		f.Write(code)
		f.Close()
		return
	}
	if verifyLicense {
		info, err := cli.GetLicenseInfo()
		if err != nil {
			panic(err)
		}
		fmt.Println(info)
		return
	}
}
