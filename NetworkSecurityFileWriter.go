package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	androidSecurityConfigFileName = "network_security_config.xml"
	iosSecurityConfigFileName     = "Info.plist"
)

type NetworkSecurityFileWriter struct{}

func (writer *NetworkSecurityFileWriter) WriteNetworkSecurityFiles(networkSecurityFiles NetworkSecurityFileBytes, outputPath string, noFile bool) {
	if noFile {
		if len(networkSecurityFiles.AndroidNetworkSecurityFile) > 0 {
			os.Stdout.Write(networkSecurityFiles.AndroidNetworkSecurityFile)
			fmt.Printf("\n\n")
		}
		if len(networkSecurityFiles.IosNetworkSecurityFile) > 0 {
			os.Stdout.Write(networkSecurityFiles.IosNetworkSecurityFile)
			fmt.Printf("\n\n")
		}
	} else {
		if len(networkSecurityFiles.AndroidNetworkSecurityFile) > 0 {
			var file *os.File
			path := filepath.Join(outputPath, androidSecurityConfigFileName)
			file, err := os.Create(path)
			if err != nil {
				fmt.Printf("Error creating file %v", err)
			}
			defer file.Close()

			file.Write(networkSecurityFiles.AndroidNetworkSecurityFile)
		}
		if len(networkSecurityFiles.IosNetworkSecurityFile) > 0 {
			var file *os.File
			path := filepath.Join(outputPath, iosSecurityConfigFileName)
			file, err := os.Create(path)
			if err != nil {
				fmt.Printf("Error creating file %v", err)
			}
			defer file.Close()

			file.Write(networkSecurityFiles.IosNetworkSecurityFile)
		}
	}
}
