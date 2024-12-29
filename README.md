# Bayaan Logger

Bayaan Logger is a highly customizable, asynchronous logging library for Go. It supports multiple log levels, structured fields, and dynamic configuration, making it suitable for modern applications.

## Features

- Multiple log levels: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC.
- Color-coded output for better visibility.
- Supports multiple outputs (e.g., stdout, files).
- Thread-safe and asynchronous logging.
- Structured fields for enhanced context.
- Customizable time formats.

## Installation

```bash
go get github.com/ahmedsat/bayaan
```

## Usage

### Basic Usage

```go
package main

import (
	"github.com/ahmedsat/bayaan"
)

func main() {
	// Setup the logger with default settings
	bayaan.Setup()

	// Log messages
	bayaan.Info("Application started", nil)
	bayaan.Warn("Low disk space", bayaan.Fields{"disk": "80%"})
	bayaan.Error("Failed to connect to database", bayaan.Fields{"error": "timeout"})

	// Close the logger when done
	bayaan.Close()
}
```

### Custom Configuration

```go
package main

import (
	"os"
	"github.com/ahmedsat/bayaan"
)

func main() {
	// Setup the logger with custom settings
	bayaan.Setup(
		bayaan.WithLevel(bayaan.LoggerLevelDebug),
		bayaan.WithTimeFormat("2006-01-02 15:04:05"),
		bayaan.WithOutput(os.Stdout, false, true),
		bayaan.WithFields(bayaan.Fields{
			"app": "Bayaan Example",
			"env": "development",
		}),
	)

	// Log messages
	bayaan.Debug("Debugging application", nil)
	bayaan.Info("Application running", nil)

	// Close the logger when done
	bayaan.Close()
}
```

### Using Structured Fields

```go
bayaan.Info("User logged in", bayaan.Fields{
	"user_id": 123,
	"role": "admin",
})
```

## Log Levels

Bayaan Logger supports the following log levels:

- `TRACE`: Detailed information for debugging.
- `DEBUG`: Debug-level messages.
- `INFO`: General operational information.
- `WARN`: Indications of potential issues.
- `ERROR`: Errors that need attention.
- `FATAL`: Critical errors causing program termination.
- `PANIC`: Severe errors causing panic.

## Features

- Multiple log levels: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC.
- Color-coded output for better visibility.
- Supports multiple outputs (e.g., stdout, files).
- Thread-safe and asynchronous logging.
- Structured fields for enhanced context.
- Customizable time formats.


## Contributing

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-name`).
3. Commit your changes (`git commit -m 'Add feature-name'`).
4. Push to the branch (`git push origin feature-name`).
5. Open a pull request.

## License

This project is licensed under the `GNU General Public Version 2` License. See the [LICENSE](LICENSE) file for details.
