## pocketix-go

A bachelor thesis project that deals with the backend processing of visual programming language made at BUT FIT.

# Overview

**Pocketix-go** is an interpret for the Pocketix v2 ecosystem. It enables the evaluation, interpretation, and execution of programs written in Pocketix, transforming them into actionable commands for various devices and systems. The project provides a modular and extensible backend architecture, making it easy to integrate new program blocks, device types, and operations.

# Features

- **Visual Programming Language Support**: Execute programs created in the Pocketix v2 visual programming environment.
- **Modular Design**: Organized into modules for parsing, statements, variables, and more.
- **Extensive Testing**: Includes unit tests for key components to ensure reliability.
- **Customizable**: Easily extendable for new program blocks, device types, and operations.


# Usage

To get started, create instances of the following:

- **models.VariableStore**: Handles user variables in the program
- **models.ProcedureStore**: Handles user procedures in the program
- **models.ReferencedValueStore**: Handles referenced values to the devices parameters
- **ast**: slices of program statements
- **statements.Collector**: collector for the program statements

Then to parse the program using:

```go
err = parser.Parse(
    modifiedData, 
    variableStore, 
    procedureStore, 
    referencedValueStore, 
    collector
)
```

To execute the commands cycle through each statement in ast:

```go
var interpretInvocationsToSend []models.SDCommandInvocation
for _, block := range ast {
    if _, err := block.Execute(variableStore, referencedValueStore, collector.DeviceCommands, func(deviceCommand models.SDCommandInvocation) {
        interpretInvocationsToSend = append(interpretInvocationsToSend, deviceCommand)
    }); err != nil {
        log.Fatalln(err)
    }
}
```

where the last argument in Execute function is a callback to collect device commands.

# Run

To run the application, run the following command in root directory:
```
go run src/main.go --path=<path_to_json_program>
```

Tests are located in /tests/ directory, to execute them run:
```
go test -v ./tests/...
```
To execute tests with coverage, you can also run tests through Makefile:
```
make test
```

# License

This project is licensed under the MIT License.
See the [LICENSE](LICENSE) file for more information.