#Tacklebox

Tacklebox is a tool that allows you to copy common project files to directories specified by a configuration.
An example use of this is for editor configurations, git hooks, and other items that aren't typically version controlled.


##Config
A sample config, with explanations is as follows:
```
{
    "GlobalTemplates": { // Global templates available by string
        "PreCommit": {
            "Location": ".git/hooks/pre-commit", // Location relative to the project to copy the file to
            "Settings": {}// Settings for the global template
        }
    },
    "Projects": [
        {
            "Name": "Testing",
            "Location": "/Projects/testing", // Location of project
            "TemplateSettings": [ // Configured templates to be copied to the project
                {
                    "Name": "VSCodeConfiguration", // Name of the template to configure
                    "Location": ".vscode", // Location relative to the project to copy the file to
                    "Settings": { // Settings for the template
                        "name": "Test Project"
                    }
                }
            ],
            "Globals": ["PreCommit"] // All of the globals to copy for this project
        }
    ],
    "Templates": [
        {
            "Name": "PreCommitHook", // Name to identify the template by
            "Location": "/Users/home/.tacklebox/templates/pre-commit" // Location of the file to copy
        },
        {
            "Name": "VSCodeConfiguration",
            "Location": "/Users/home/.tacklebox/templates/.vscode"
        }
    ]
}
```

##Commands

###Sync
Usage: `tacklebox sync` or `tacklebox s`
Using the current configuration file, copies all files configured to projects. Given the [above configuration file](#config),
the sync command would copy the `/Users/home/.tacklebox/templates/pre-commit` file to `/Projects/testing/.git/hooks/pre-commit`,
and would copy the the `/Users/home/.tacklebox/templates/.vscode` file to `/Projects/testing/.vscode`

###ReadDir
Usage: `tacklebox readdir <DirToRead>` or `tacklebox r <DirToRead>`
Given a `DirToRead`, adds all directories inside `DirToRead` to the config file as a empty project by the name of the directory.