# icaljson

```go
import "github.com/beyondcivic/icaljson/pkg/icaljson"
```

## Index

- [func IsICalFile\(filePath string\) bool](<#IsICalFile>)
- [func ValidateOutputPath\(outputPath string\) error](<#ValidateOutputPath>)
- [type AppError](<#AppError>)
  - [func \(e AppError\) Error\(\) string](<#AppError.Error>)
- [type Calendar](<#Calendar>)
  - [func Generate\(icsPath string, outputPath string\) \(\*Calendar, error\)](<#Generate>)
- [type Event](<#Event>)


<a name="IsICalFile"></a>
## func [IsICalFile](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/utils.go#L11>)

```go
func IsICalFile(filePath string) bool
```

IsCSVFile checks if a file appears to be a CSV file based on extension

<a name="ValidateOutputPath"></a>
## func [ValidateOutputPath](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/utils.go#L17>)

```go
func ValidateOutputPath(outputPath string) error
```

ValidateOutputPath validates if the given path is a valid file path

<a name="AppError"></a>
## type [AppError](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/error.go#L5-L10>)



```go
type AppError struct {
    // Message to show the user.
    Message string
    // Value to include with message
    Value any
}
```

<a name="AppError.Error"></a>
### func \(AppError\) [Error](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/error.go#L12>)

```go
func (e AppError) Error() string
```



<a name="Calendar"></a>
## type [Calendar](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/structs.go#L3-L7>)



```go
type Calendar struct {
    ProdID  string  `json:"prodid,omitempty"`
    Version string  `json:"version,omitempty"`
    Events  []Event `json:"events,omitempty"`
}
```

<a name="Generate"></a>
### func [Generate](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/core.go#L10>)

```go
func Generate(icsPath string, outputPath string) (*Calendar, error)
```

Generate generates JSON file from a ICS file with automatic type inference.

<a name="Event"></a>
## type [Event](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/structs.go#L9-L19>)



```go
type Event struct {
    UID         string   `json:"uid,omitempty"`
    Summary     string   `json:"summary,omitempty"`
    Description string   `json:"description,omitempty"`
    Start       string   `json:"start,omitempty"`
    End         string   `json:"end,omitempty"`
    Location    string   `json:"location,omitempty"`
    URL         string   `json:"url,omitempty"`
    Status      string   `json:"status,omitempty"`
    Categories  []string `json:"categories,omitempty"`
}
```