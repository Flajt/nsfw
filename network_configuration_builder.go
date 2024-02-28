package main

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"fmt"
	"text/template"
	"time"
)

const (
	YYYYMMDD = "2006-01-02" // No freaking clue why it has to be this date, 06 for year, 01 for month, 02 for day?
)

//go:embed tmpl/NSAppTransportSecurity.tmpl
var nsAppTransportSecurityTmpl string //Not sure if that is the right way to place this

type NetworkSecurityFileBuilder struct{}

type AndroidDomainConfigWrapper struct {
	XMLName xml.Name      `xml:"domain-config"`
	Domain  AndroidDomain `xml:"domain"`
	PinSet  AndroidPinSet `xml:"pin-set"`
}
type AndroidNetworkSecurityWrapper struct {
	XMLName xml.Name                     `xml:"network-security-config"`
	Config  []AndroidDomainConfigWrapper `xml:"domain-config"`
}
type AndroidDomain struct {
	Domain     string `xml:",chardata"`
	IncludeSub bool   `xml:"includeSubdomains,attr"`
}

type AndroidPinSet struct {
	Expiration string     `xml:"expiration,attr"`
	Pin        AndroidPin `xml:"pin"`
}

type AndroidPin struct {
	DigestType string `xml:"digest,attr"`
	Digest     string `xml:",chardata"`
}

// / Contains all required data for network security configuration, currently only support the last part of the certificate chain
type CertificateInfo struct {
	FingerPrint string
	Domain      string
	expiryDate  time.Time
}

type NetworkSecurityFileBytes struct {
	AndroidNetworkSecurityFile []byte
	IosNetworkSecurityFile     []byte
}

// / Wraps ios and android builder logic
func (builder *NetworkSecurityFileBuilder) GenerateCertificateInfo(CertificateInfo []CertificateInfo, platforms []string) NetworkSecurityFileBytes {
	var androidSecurityConfig bytes.Buffer
	var iosSecurityConfig bytes.Buffer
	for _, platform := range platforms {
		switch platform {
		case "android":
			androidSecurityConfig = builder.generateAndroidCertificateInfo(CertificateInfo)

		case "ios":
			iosSecurityConfig = builder.generateIosCertificateInfo(CertificateInfo)
		}
	}
	return NetworkSecurityFileBytes{
		AndroidNetworkSecurityFile: androidSecurityConfig.Bytes(),
		IosNetworkSecurityFile:     iosSecurityConfig.Bytes(),
	}
}

func (builder NetworkSecurityFileBuilder) generateAndroidCertificateInfo(CertificateInfo []CertificateInfo) bytes.Buffer {
	wrapper := AndroidNetworkSecurityWrapper{}
	for _, cfinfo := range CertificateInfo {
		date := cfinfo.expiryDate.Format(YYYYMMDD)
		domain := AndroidDomain{
			Domain:     cfinfo.Domain,
			IncludeSub: true,
		}
		pin := AndroidPin{
			DigestType: "SHA-256",
			Digest:     cfinfo.FingerPrint,
		}
		pinSet := AndroidPinSet{
			Expiration: date,
			Pin:        pin,
		}
		domainConfig := AndroidDomainConfigWrapper{
			Domain: domain,
			PinSet: pinSet,
		}
		// / Add the entry to the AndroidNetworkSecurityWrapper
		wrapper.Config = append(wrapper.Config, domainConfig)
	}
	xmlBytes, err := xml.MarshalIndent(wrapper, " ", "	")
	xmlBytes = []byte(xml.Header + string(xmlBytes))
	if err != nil {
		fmt.Printf("Error marshaling xml %v", err)
	}
	var buffer bytes.Buffer
	buffer.Write(xmlBytes)
	return buffer
}

func (builder NetworkSecurityFileBuilder) generateIosCertificateInfo(CertificateInfo []CertificateInfo) bytes.Buffer {
	const templateName string = "NSAppTransportSecurity.tmpl"

	tmpl, err := template.New(templateName).Parse(nsAppTransportSecurityTmpl)
	if err != nil {
		fmt.Printf("Error parsing template %v", err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, CertificateInfo)
	if err != nil {
		fmt.Printf("Error executing template %v", err)
	}
	return buffer

}
