package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	var websitesRaw string
	var outputPath string
	var platformsRaw string
	var noFile bool
	nsfw := NetworkSecurityFileWriter{}

	// add flags: -websites for website, -output for output location, -platforms for android or ios
	flag.StringVar(&websitesRaw, "websites", "", "Websites to add to network security configuration files, COMMA separated!")
	flag.StringVar(&outputPath, "output", ".", "Output location for network security configuration files!")
	flag.StringVar(&platformsRaw, "platforms", "android,ios", "Platforms to generate network security configuration files for,android & ios are options, comma separated!")
	flag.BoolVar(&noFile, "no-file", false, "Prints the network security configuration files to stdout instead of writing to a file")
	flag.Parse()

	if websitesRaw == "" {
		fmt.Println("Please provide websites to add to network security configuration files!")
	}
	if outputPath == "" {
		fmt.Println("Please provide output location for network security configuration files!")
	}
	if platformsRaw == "" {
		fmt.Println("Please provide platforms to generate network security configuration files for!")
	}
	platforms := strings.Split(platformsRaw, ",")
	websites := strings.Split(websitesRaw, ",")
	client := http.Client{}
	certificates := make([]CertificateInfo, len(websites))
	if len(websites) > 0 {
		fmt.Println("Gathering TLS certificates for:", websites)
		fmt.Println("")
		for i, site := range websites { //TODO: Consider moving this to a function
			req, reqErr := http.NewRequest(http.MethodOptions, site, nil)
			if reqErr != nil {
				fmt.Printf("Error creating request %v", reqErr)
			}
			resp, err := client.Do(req)
			if resp.TLS != nil {
				peerCertificate := resp.TLS.PeerCertificates[0]
				fingerPrintRaw := sha256.Sum256(peerCertificate.RawSubjectPublicKeyInfo)
				base64FingerPrint := base64.StdEncoding.EncodeToString(fingerPrintRaw[:])
				certificates[i] = CertificateInfo{
					FingerPrint: base64FingerPrint,
					Domain:      resp.TLS.ServerName,
					expiryDate:  resp.TLS.PeerCertificates[0].NotAfter,
				}

			}
			defer resp.Body.Close()
			if err != nil {
				fmt.Println("Error fetching website", err)
			}

		}
		NetworkSecurityFileBuilder := NetworkSecurityFileBuilder{}
		NetworkSecurityFiles := NetworkSecurityFileBuilder.GenerateCertificateInfo(certificates, platforms)
		nsfw.WriteNetworkSecurityFiles(NetworkSecurityFiles, outputPath, noFile)

	} else {
		fmt.Println("No websites to fetch!")
	}
}
