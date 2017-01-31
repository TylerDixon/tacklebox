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
	GlobalTemplates map[string]TemplateSetting // Map of string IDs to template settings for global use
	Projects        []ProjectConfig            // Collection of projects
	Templates       []Template                 // Templates that can be configured for projects
}

type ProjectConfig struct {
	Globals          []string          // List of globals for this project
	Name             string            // Unique name of the project
	Location         string            // Location of the root of the project
	TemplateSettings []TemplateSetting // Collection of templates to render for the project and their settings
}

type TemplateSetting struct {
	Name     string                 // Name identifier of the template to render
	Location string                 // Location in the project the template should be stored
	Settings map[string]interface{} // Settings to render the template with
}

// Sync the all files declared in configuration
func (configData ConfigData) Sync() error {
	filesToSync := make(map[string][]byte)
	for _, project := range configData.Projects {
		// Render all templates configured for the template
		for _, templateSetting := range project.TemplateSettings {
			templateToRender, getTemplateError := getTemplateByName(configData.Templates, templateSetting.Name)
			if getTemplateError != nil {
				err := fmt.Errorf("Failed to find templateToRender for project %s due to error %s", project.Name, getTemplateError)
				fmt.Println(err)
				return err
			}

			renderedTemplate, renderErr := templateToRender.Render(templateSetting.Settings)
			if renderErr != nil {
				err := fmt.Errorf("Failed to sync files due to render error %s", renderErr)
				fmt.Println(err)
				return err
			}
			filesToSync[path.Join(project.Location, templateSetting.Location)] = renderedTemplate
		}

		// Render all global templates specified for the project
		for _, global := range project.Globals {
			if globalConfig, ok := configData.GlobalTemplates[global]; ok {
				templateToRender, getTemplateError := getTemplateByName(configData.Templates, globalConfig.Name)
				if getTemplateError != nil {
					err := fmt.Errorf("Failed to find templateToRender for global %s due to error %s", globalConfig.Name, getTemplateError)
					fmt.Println(err)
					return err
				}

				renderedTemplate, renderErr := templateToRender.Render(globalConfig.Settings)
				if renderErr != nil {
					err := fmt.Errorf("Failed to sync files due to render error %s", renderErr)
					fmt.Println(err)
					return err
				}
				filesToSync[path.Join(project.Location, globalConfig.Location)] = renderedTemplate

			} else {
				err := fmt.Errorf("Project %s configured with global %s, but no such global exists.", project.Name, global)
				fmt.Println(err)
				return err
			}
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

// ConfigDirs takes in a directory to read, and for each directory in it, adds a project entry to config
// TODO: Add interactive checkbox selection for adding dirs.
// TODO: Check if dir already exists in config.
func (configData *ConfigData) ConfigDirs(dirToRead string) error {
	fileInfo, readDirErr := ioutil.ReadDir(dirToRead)
	if readDirErr != nil {
		err := fmt.Errorf("Failed to ReadDir %s due to error %s", dirToRead, readDirErr)
		fmt.Println(err)
		return err
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			configData.Projects = append(configData.Projects, ProjectConfig{
				Name:             file.Name(),
				Location:         path.Join(dirToRead, file.Name()),
				TemplateSettings: []TemplateSetting{},
			})
		}
	}
	return nil
}

// Retrieve config returns a parsed config data from the tacklebox config directory in the home directory.
func RetrieveConfig() (ConfigData, error) {
	var unmarshaledConfigData ConfigData
	//TODO: Handle homedir errors
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

// Write config data to config file
func (configData ConfigData) Save() error {
	//TODO: Handle homedir errors
	homeDirectory, _ := homedir.Dir()
	configFilePath := path.Join(homeDirectory, ".tacklebox", "config.json")
	marshalledConfigData, marshalErr := json.MarshalIndent(configData, "", "    ")
	if marshalErr != nil {
		err := fmt.Errorf("Failed to marshal config data due to error %s", marshalErr)
		fmt.Print(err)
		return err
	}
	writeFileErr := ioutil.WriteFile(configFilePath, marshalledConfigData, 0666)
	if writeFileErr != nil {
		err := fmt.Errorf("Failed to write config data due to error %s", writeFileErr)
		fmt.Println(err)
		return err
	}
	fmt.Printf("Wrote config to file %s\n", configFilePath)
	return nil
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
