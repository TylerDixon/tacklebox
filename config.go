package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type ConfigData struct {
	Projects  []ProjectConfig // Collection of projects
	Templates []Template      // Templates that can be configured for projects
}

type ProjectConfig struct {
	Name             string            // Unique name of the project
	Location         string            // Location of the root of the project
	TemplateSettings []TemplateSetting // Collection of templates to render for the project and their settings
}

type TemplateSetting struct {
	Name     string                 // Name identifier of the template to render
	Location string                 // Location in the project the template should be stored TODO: Should be relative to ProjectConfig.Location
	Settings map[string]interface{} // Settings to render the template with
}

// Sync the all files declared in configuration
func (configData ConfigData) sync() error {
	filesToSync := make(map[string][]byte)
	for _, project := range configData.Projects {
		for _, templateSetting := range project.TemplateSettings {
			templateToRender, getTemplateError := getTemplateByName(configData.Templates, templateSetting.Name)
			if getTemplateError != nil {
				return fmt.Errorf("Failed to find tepmlate for project %s due to error %s", project.Name, getTemplateError)
			}
			renderedTemplate, renderErr := templateToRender.Render(templateSetting.Settings)
			if renderErr != nil {
				return fmt.Errorf("Failed to sync files due to render error %s", renderErr)
			}
			filesToSync[templateSetting.Location] = renderedTemplate
		}
	}
	for writeLocation, fileData := range filesToSync {
		writeRenderedFileError := ioutil.WriteFile(writeLocation, fileData, 0666)
		if writeRenderedFileError != nil {
			err := fmt.Errorf("Failed to write to file %s due to error %s", writeLocation, writeRenderedFileError)
			fmt.Println(err)
			return err
		} else {
			fmt.Printf("Synced file at %s\n", writeLocation)
		}
	}
	return nil
}

// Retrieve config returns a parsed config data from the tacklebox config directory in the home directory.
func RetrieveConfig() (ConfigData, error) {
	var unmarshaledConfigData ConfigData
	homeDirectory, _ := homedir.Dir()
	mkDirErr := os.Mkdir(path.Join(homeDirectory, ".tacklebox"), 0777)
	if !os.IsExist(mkDirErr) {
		fmt.Printf("Create tacklebox config directory error: %s", mkDirErr)
		return unmarshaledConfigData, mkDirErr
	}
	configFilePath := path.Join(homeDirectory, ".tacklebox", "config.json")

	// Create or retrieve config data
	var encodedConfigData []byte
	var fileCreationOrReadError error
	if _, err := os.Stat(configFilePath); err == nil {
		encodedConfigData, fileCreationOrReadError = ioutil.ReadFile(configFilePath)
	} else {
		encodedConfigData, fileCreationOrReadError = InitializeConfigFile(configFilePath)
	}
	if fileCreationOrReadError != nil {
		return unmarshaledConfigData, fileCreationOrReadError
	}

	unmarshalConfigError := json.Unmarshal(encodedConfigData, &unmarshaledConfigData)

	if unmarshalConfigError != nil {
		fmt.Printf("Failed to unmarshal config data with error %s", unmarshalConfigError)
		return unmarshaledConfigData, unmarshalConfigError
	}
	return unmarshaledConfigData, nil
}

// InitializeConfigFile creates a default configuration file for use by tacklebox.
func InitializeConfigFile(configFilePath string) ([]byte, error) {
	file, createFileErr := os.Create(configFilePath)
	if createFileErr != nil {
		log.Fatalf("Failed to create config file at %s due to error: %s", configFilePath, createFileErr)
		return nil, createFileErr
	}
	defaultConfigFile := []byte(`{
		"Templates": []
	}`)
	_, writeFileErr := file.Write(defaultConfigFile)
	if writeFileErr != nil {
		log.Fatalf("Failed to create config file at %s due to error: %s", configFilePath, writeFileErr)
		return nil, writeFileErr
	}
	return defaultConfigFile, nil
}
