## **Golang Command line**

<div align="center">
    <a href="#">  <img height="150" src="https://raw.githubusercontent.com/betandr/gophers/master/Gopher.png" alt="" /></a>
</div><br>

## Description
Gcli is a command line made in golang to help you create crud in your api project more easily. It is also possible to create a new project with the necessary structure to start your api project.

## Templates

| Template | Description | Available |
| --- | --- | --- |
| [gcli-lite-template](https://github.com/vinirossado/gcli-lite-template) | A lite template with only the necessary to start a project with database connection and a simple handler. | ❌ |
| [gcli-basic-template](https://github.com/vinirossado/gcli-basic-template ) | A basic template with the necessary to start a project with gorm and postgres connection. | ❌ |
| [gcli-advanced-template](https://github.com/vinirossado/gcli-advanced-template) | A advanced template with the necessary to start a expansive project with gorm, postgres, jwt, cors, logger, wire and more. | ✅ |


## Install and Run.
1. Clone the repository : `git clone https://github.com/vinirossado/gcli`.
2. Install [GO](https://go.dev/) to run.
3. Install [Visual Studio Code](https://code.visualstudio.com/) to edit or use your favorite editor.
4. Open your Terminal and run ```go mod tidy``` in folder to install dependencies.
5. Run ```go run main.go```.
6. Run ```go build main.go``` to build the project.
7. Install the cli with ```go install```.
8. Run ```gcli``` to see the commands.


# Build with Makefile (Optional).
1. Install [Make](https://www.gnu.org/software/make/) to run.
2. Run ```make build linux``` to build the project for linux.
3. Run ```make build windows``` to build the project for windows.
4. Run ```make build darwin``` to build the project for darwin.


## Commands
- ```gcli``` to see the commands.
- ```gcli create``` Create a new handler/service/repository/model.
- ```gcli new``` Create a new api project template.
- ```gcli version``` Show the cli version.
- ```gcli upgrade``` Upgrade the cli to the latest version.
- ```gcli help <command>``` Help about any command.
- ```gcli wire cmd/server``` [wire.go path] Wire is a code generation tool that automates connecting components using dependency injection.

Gopher artwork by [Beth Anderson](https://github.com/betandr/gophers).

## Contributing

1. [Fork the repository](https://github.com/vinirossado/gcli/fork)!
2. Clone your fork.
3. Create your feature branch: `git checkout -b my-new-feature`
4. Commit your changes: `git commit -am 'Add some feature'`
5. Push to the branch: `git push origin my-new-feature`
6. Submit a pull request :D