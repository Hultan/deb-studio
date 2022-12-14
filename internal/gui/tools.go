package gui

import (
	"fmt"
	"io"
	"os"

	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/deb-studio/internal/common"
)

// getApplicationName: Get the application name, version and copyright
func getApplicationName() string {
	log.Trace.Println("Entering getApplicationName...")
	defer log.Trace.Println("Exiting getApplicationName...")

	return fmt.Sprintf(
		"%s %s - %s",
		common.ApplicationTitle,
		common.ApplicationVersion,
		common.ApplicationCopyRight,
	)
}

func getProjectStatus() projectStatus {
	log.Trace.Println("Entering getProjectStatus...")
	defer log.Trace.Println("Exiting getProjectStatus...")

	if project == nil {
		return projectStatusNoProjectOpened
	} else if project.CurrentPackage == nil {
		return projectStatusNoPackageSelected
	} else if !project.IsWorkingWithLatestVersion() {
		return projectStatusNotLatestVersion
	}
	return projectStatusLatestVersion
}

func readTextFile(path string) (string, error) {
	log.Trace.Println("Entering readTextFile...")
	defer log.Trace.Println("Exiting readTextFile...")

	f, err := os.Open(path)
	if err != nil {
		log.Error.Printf("failed to open file '%s': %s\n", path, err)
		return "", err
	}
	text, err := io.ReadAll(f)
	if err != nil {
		log.Error.Printf("failed to read all from file '%s': %s\n", path, err)
		return "", err
	}

	return string(text), nil
}

func writeTextFile(path, content string) error {
	log.Trace.Println("Entering writeTextFile...")
	defer log.Trace.Println("Exiting writeTextFile...")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Error.Printf("failed to open file '%s': %s\n", path, err)
		return err
	}
	n, err := io.WriteString(f, content)
	if err != nil {
		log.Error.Printf("failed to write to file '%s': %s\n", path, err)
		return err
	}
	if n != len(content) {
		log.Error.Printf("failed to write complete file '%s': %s\n", path, err)
		return err
	}

	return nil
}

func setTextViewText(tv *gtk.TextView, text string) {
	log.Trace.Println("Entering setTextViewText...")
	defer log.Trace.Println("Exiting setTextViewText...")

	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Error.Printf("failed to get text view buffer: %s\n", err)
		return
	}
	buffer.SetText(text)
}

func getTextViewText(tv *gtk.TextView) (string, error) {
	log.Trace.Println("Entering getTextViewText...")
	defer log.Trace.Println("Exiting getTextViewText...")

	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Error.Printf("failed to get text view buffer: %s\n", err)
		return "", err
	}
	text, err := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if err != nil {
		log.Error.Printf("failed to get text view text: %s\n", err)
		return "", err
	}
	return text, nil
}

func getEntryText(e *gtk.Entry) (string, error) {
	log.Trace.Println("Entering getEntryText...")
	defer log.Trace.Println("Exiting getEntryText...")

	text, err := e.GetText()
	if err != nil {
		log.Error.Printf("failed to get entry text: %s\n", err)
		return "", err
	}
	return text, nil
}

func getEntryName(e *gtk.Entry) (string, error) {
	log.Trace.Println("Entering getEntryName...")
	defer log.Trace.Println("Exiting getEntryName...")

	name, err := e.GetName()
	if err != nil {
		log.Error.Printf("failed to get entry name: %s\n", err)
		return "", err
	}
	return name, nil
}
