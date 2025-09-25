# AGENTS.md - Guide for Agentic Coding Agents

## Build Commands

- Build the project: `go build`
- Install dependencies: `go mod tidy`

## Test Commands

- Run all tests: `go test -v`
- Run a single test: `go test -run TestName`
- Run tests with coverage: `go test -coverprofile=coverage.out && go tool cover
   -html=coverage.out -o coverage.html`

## Lint Commands

- Format code: `gofmt -w .`
- Vet code: `go vet ./...`
- Check formatting: `gofmt -d .`

## Application Context

- **Overview**: This is a JSON API for retrieving medicaments data from parsed external sources. It provides endpoints for searching, listing, and fetching medicaments and generics.
- **HTML Usage**: HTML files (e.g., `html/index.html`) are used only for static documentation and are served with caching headers. No user inputs are processed or rendered in HTML, reducing XSS risks.
- **Input Sanitization**: User inputs (e.g., search terms) are validated for length and content (alphanumeric + spaces), and sanitized for regex use with `regexp.QuoteMeta`. JSON marshaling handles escaping automatically for API responses.
- **Security Notes**: Focus on API security; no database interactions, so SQL injection is not a concern. Rate limiting and middleware protect against abuse.

## Code Style Guidelines

- **Imports**: Group standard library, third-party, then local packages.
- **Formatting**: Use `gofmt` for consistent indentation and spacing.
- **Naming**: Exported identifiers use CamelCase; unexported use camelCase.
- **Types**: Use explicit types; prefer structs for data.
- **Error Handling**: Use `log.Fatal` for critical errors; `log.Printf` for warnings.
- **Concurrency**: Use goroutines and channels for parallel processing.
- **Comments**: Add comments for exported functions and complex logic.
