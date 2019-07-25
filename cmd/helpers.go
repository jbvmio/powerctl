package cmd

import (
	"io/ioutil"
	"os"
)

// FileExists checks for the existence of the file indicated by filename and returns true if it exists.
func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func copyKubeConfig(kubeConfigPath, pctlPath string) error {
	kubeConfig, err := readFile(kubeConfigPath)
	if err != nil {
		return err
	}
	return writeFile(pctlPath, kubeConfig)
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func writeFile(path string, file []byte) error {
	return ioutil.WriteFile(path, file, 0644)
}
