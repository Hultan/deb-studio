package engine

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type descriptorType int

const (
	projectDescriptor descriptorType = iota
	versionDescriptor
	architectureDescriptor
)

func haveDescriptor(descType descriptorType, descPath string) bool {
	log.Trace.Println("Entering haveDescriptor...")
	defer log.Trace.Println("Exiting haveDescriptor...")

	descName := getDescriptorFileName(descType)
	filePath := path.Join(descPath, descName)

	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.IsDir() == false
}

func readDescriptor(descType descriptorType, descPath string) (string, error) {
	log.Trace.Println("Entering readDescriptor...")
	defer log.Trace.Println("Exiting readDescriptor...")

	descName := getDescriptorFileName(descType)
	descPrefix := getDescriptorPrefix(descType)
	filePath := path.Join(descPath, descName)

	text, err := readAllText(filePath)
	if err != nil {
		log.Error.Printf("Failed to read all text of path '%s': %s\n", filePath, err)
		return "", err
	}

	name, err := getFirstLine(text, descPrefix, "\n")
	if err != nil {
		log.Error.Printf("Failed to get first line of path '%s': %s\n", filePath, err)
		return "", err
	}

	return strings.Trim(name, " \t"), nil
}

func writeDescriptor(descType descriptorType, descPath, name string) error {
	log.Trace.Println("Entering writeDescriptor...")
	defer log.Trace.Println("Exiting writeDescriptor...")

	descName := getDescriptorFileName(descType)
	descPrefix := getDescriptorPrefix(descType)
	filePath := path.Join(descPath, descName)

	content := fmt.Sprintf("%s=%s", descPrefix, name)
	err := createTextFile(filePath, content)
	if err != nil {
		log.Error.Printf("Failed to write descriptor to path '%s': %s\n", filePath, err)
		return err
	}
	return nil
}

func getDescriptorFileName(descType descriptorType) string {
	switch descType {
	case projectDescriptor:
		return ".project"
	case versionDescriptor:
		return ".version"
	case architectureDescriptor:
		return ".architecture"
	default:
		return ""
	}
}

func getDescriptorPrefix(descType descriptorType) string {
	switch descType {
	case projectDescriptor:
		return "PROJECT"
	case versionDescriptor:
		return "VERSION"
	case architectureDescriptor:
		return "ARCHITECTURE"
	default:
		return ""
	}
}
