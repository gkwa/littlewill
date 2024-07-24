# littlewill

Purpose:

littlewill is a tool for cleaning up Markdown links by removing extra spaces within the link text.

## Example Usage

```bash
littlewill input.md > output.md
```

## Install littlewill

On macOS/Linux:
```bash
brew install gkwa/homebrew-tools/littlewill
```

On Windows:
```powershell
TBD
```

## Running Tests

To run the tests for littlewill, follow these steps:

1. Ensure you have Go installed on your system.
2. Navigate to the project root directory.
3. Run the following command:

```bash
go test ./...
```

This will run all tests in the project, including the ones in the `core` package.

To run tests with verbose output:

```bash
go test -v ./...
```

To run a specific test:

```bash
go test -v ./core -run TestCleanupMarkdownLinks
```

Replace `TestCleanupMarkdownLinks` with the name of the specific test you want to run.

## Development

The main functionality is implemented in the `core/cleanup.go` file. Tests are located in `core/cleanup_test.go`.

To contribute or modify the code, make sure to run the tests after any changes to ensure everything is working as expected.
