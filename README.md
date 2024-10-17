# Go Coverage Plus
Optimise the coverage report in go text format with source code.
Also, can write cobertura format. Including complexity metrics and branch coverage.

## Installation
Simple install it by `go install`:
```
go install https://github.com/Fabianexe/gocoverageplus@latest
```

## Usage
`gocoverageplus` can be run without any arguments (fallbacking to defaults).
However, it needs a config file in json format (See config section).
You can further specify the config path, a coverage input path, and an output path.

### Flags
* `-h` or `--help` to get a help message
* `-c` or `--config` to specify the config file (default is `.cov.json`)
* `-o` or `--output` to specify the output file (default is `coverage.xml`)
* `-i` or `--input` to specify the input file (default is `coverage.cov"`)
* `-v` or `--verbose` to get more output. Can be used multiple times to increase the verbosity

### Config
The config file is a json file with the following structure:
```json
{
  "OutputFormat": "cobertura",
  "SourcePath": "./",
  "Cleaner": {
    "ErrorIf": true,
    "NoneCodeLines": true,
    "Generated": true,
    "CustomIf": [
      "debug"
    ]
  },
  "Complexity": {
    "Active": true,
    "Type": "cognitive"
  }
}
```
As output formats, you can choose between `cobertura` and `textfmt`. The first is described in the cobertura section and 
the second is the default go coverage format.
Complexity only apply for the `cobertura` format. The complexity type can be either `cognitive` or `cyclomatic`. 
The difference between these metrics is described in the cobertura section.

## The accuracy of `go test -coverprofile`
The `go test -coverprofile` command is a great tool to get coverage information about your project.
However, it measures the coverage on a bock level. This means that if you function contains empty lines, only comments,
or lines with only a closing bracket, they will be counted in line metrics.

This project tries to solve this problem by using the `go/ast` package to determine the actual lines of code from the source.

Another result from this is that branches on a line level can be determined. If a line contains an `if` statement,
with multiple conditions, it is still one block for the coverage profile. There are projects that try to solve this problem
for example [gobco](https://github.com/rillig/gobco). However, they for the moment not compatible with the Jenkins coverage plugin.
Thus, we add branch coverage on method and file level. Where such multi condition statements are counted as one branche.

## Source code Filter
There are parts of the source code that may not be included in the coverage report.
At the moment, the following parts can be excluded:
* Generated files
    * Files that fellows [this convention](https://go.dev/s/generatedcode) are excluded
* None code lines
    * Empty lines
    * Lines that only contain a comment
    * Lines that only contain a closing bracket
* Error ifs
    * If statements that only contain an error check (`if err != nil`) with only a return in the body are excluded
* Custom ifs
    * If statements that only contain one bool with a name given in the config list (`if debug`) with no else part are excluded

You can activate these filters by using the corresponding config values.

# Cobertura Format
The cobertura format is a widely used format for coverage reports. It is supported by many tools like Jenkins.
It is an XML format that contains the coverage information for each file and package. 
Besides the coverage information, it also contains the complexity metrics for each function.
The format is described  [here](https://github.com/cobertura/cobertura/blob/master/cobertura/src/site/htdocs/xml/coverage-04.dtd).
## Cyclomatic Complexity vs Cognitive Complexity

Cyclomatic Complexity and Cognitive Complexity are both software metrics used to measure the complexity of a program. They are used to determine the quality of code and identify areas that might need refactoring. However, they approach the measurement of complexity from different perspectives.

### Cyclomatic Complexity

Cyclomatic Complexity, introduced by Thomas McCabe in 1976, is a quantitative measure of the number of linearly independent paths through a program's source code. It is computed using the control flow graph of the program. The cyclomatic complexity of a section of source code is the count of the number of linearly independent paths through the source code. It is computed as:

```
Cyclomatic Complexity = Edges - Nodes + 2*Connected Components
```

Cyclomatic Complexity is primarily used to evaluate the complexity and understandability of a program, and it can also give an idea of the number of test cases needed to achieve full branch coverage.

### Cognitive Complexity

Cognitive Complexity, introduced by SonarSource, is a measure that focuses on how difficult the code is to understand by a human reader. It considers things like the level of nesting, the number of break or continue statements, the number of conditions in a decision point, and the use of language structures that unnecessarily increase complexity.

Cognitive Complexity aims to produce a measurement that will correlate more closely with a developer's experience of a code base, making it easier to identify problematic areas of code that need refactoring.

### Summary

In summary, while Cyclomatic Complexity is a measure of the structural complexity of a program, Cognitive Complexity is a measure of how difficult a program is to understand by a human reader. 
Both are useful, but they serve different purposes and can lead to different conclusions about the code's quality.

## Others
So far we are aware about two other projects that do something similar:
* [gocov-xml](https://github.com/AlekSi/gocov-xml)
* [gocover-cobertura](https://github.com/boumenot/gocover-cobertura)

However, both of them focus on the coverage part and take over a big downsides of the `go test -coverprofile` command.
Further this project adds complexity metrics, more options to determine coverage, and branch coverage.