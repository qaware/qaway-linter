# QAway Linter

The QAway linter for [golangci-lint](https://golangci-lint.run/) enforces coding best practices, especially in terms of
documentation. It allows you to define rules for your codebase with maximum flexibility as different parts of your code
may require different rules.

The linter is not yet integrated directly into [go-langci-lint](https://golangci-lint.run/), but can be used as
a [custom plugin](https://golangci-lint.run/plugins/module-plugins/) (see below).

## Usage

1. Create a file called `.custom-gcl.yml` in your projects root directory with the following content:

```yaml
version: v1.62.2 # TODO: update to latest version (see https://github.com/golangci/golangci-lint/releases/tag/v1.62.2)
plugins:
  - module: 'github.com/qaware/qaway-linter'
    import: 'github.com/qaware/qaway-linter'
    version: v0.0.1 # TODO: use latest version from GitHub 

```

2. Execute `golangci-lint custom` in the same directory to build a custom version of golangci-lint with the
   `qaway-linter` plugin.

3. Extend the configuration in your `.golangci.yml`. Customize the rules according to your codebase. Note that more
   concrete packages override the configuration of more general packages .

```yaml
linters:
  enable:
    # add to existing list if linters.disable-all is set to true
    - qawaylinter

linter-settings:
  custom:
    qawaylinter:
      type: "module"
      description: "Checks for appropriate documentation in code"
      settings:
        rules:
          - packages: [ "github.com/myorg/myrepo" ]
            # This rule demonstrates all available configuration options
            # If a parameter is not set, it is not enforced.
            functions:
              filters:
                # Apply parameters only to functions with at least 10 lines of code
                minLinesOfCode: 10
              params:
                # A method must have at least 10% of comments (headline + inline) compared to its lines of code
                minCommentDensity: 0.1
                # A headline comment is required for every method
                requireHeadlineComment: true
                # Trivial comments (similarity to method name) are not allowed. 
                # The threshold indicates the similarity to the method name.
                # A higher threshold indicates a higher similarity, resulting in less warnings.
                trivialCommentThreshold: 0.5
                # Amount of logging statements compared to lines of code. 
                minLoggingDensity: 0.0
            interfaces:
              params:
                # A headline comment is required for every interface
                requireHeadlineComment: true
                # A comment is required for every method in an interface
                requireMethodComment: true
            structs:
              params:
                # A headline comment is required for every struct
                requireHeadlineComment: true
                # A comment is required for every field in a struct
                requireFieldComment: false
          - packages: [ "github.com/myorg/myrepo/subpkg" ] # rules for subpackage override super packages
            functions:
              filters:
                minLinesOfCode: 20
              params:
                trivialCommentThreshold: 0.5
                minLoggingDensity: 0.1
```

4. Execute the custom version by running `./custom-gcl run` in your project's root directory.

## Exclusions

Add `// nolint:qawaylinter` to the line you want to exclude from the linter. It is not possible to disable individual
rules from the configuration.

## Internal architecture

A generic `Rule` interface is defined in [rule.go](rule.go). Each rule implements this interface and is responsible for
checking a specific aspect of the code. The `Rule` interface defines three methods:

* `isApplicable`: Checks if the rule is applicable to the given code element (node).
* `Analyse`: Analyzes the code element and returns a list of findings.
* `Apply`: Validates if the results from the analysis violate rules from the configuration and reports errors.

Each rule is called in the [analyser](analyser.go) for each code element. The analyser is responsible for traversing the
code and calling the rules. New rules must be added to the analyser to be executed.