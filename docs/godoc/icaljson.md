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
- [type Geolocation](<#Geolocation>)


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
## type [Calendar](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/structs.go#L4-L15>)

Calendar represents a VCALENDAR component according to RFC 5545

```go
type Calendar struct {
    // Required properties
    ProdID  string `json:"prodid,omitempty"`  // Product Identifier
    Version string `json:"version,omitempty"` // iCalendar Version (should be 2.0)

    // Optional calendar properties
    CalScale string `json:"calscale,omitempty"` // Calendar scale (e.g., GREGORIAN)
    Method   string `json:"method,omitempty"`   // iTIP method (e.g., REQUEST, PUBLISH)

    // Components
    Events []Event `json:"events,omitempty"`
}
```

<a name="Generate"></a>
### func [Generate](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/core.go#L15>)

```go
func Generate(icsPath string, outputPath string) (*Calendar, error)
```

Generate generates JSON file from a ICS file with automatic type inference.

<a name="Event"></a>
## type [Event](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/structs.go#L23-L70>)

Event represents a VEVENT component according to RFC 5545

```go
type Event struct {
    // Required properties (in most contexts)
    UID string `json:"uid,omitempty"` // Unique identifier

    // Date/Time properties
    Start    string `json:"start,omitempty"`    // DTSTART - Start date/time
    End      string `json:"end,omitempty"`      // DTEND - End date/time
    Duration string `json:"duration,omitempty"` // DURATION - Alternative to DTEND

    // Core descriptive properties
    Summary     string `json:"summary,omitempty"`     // Brief description/title
    Description string `json:"description,omitempty"` // Full description
    Location    string `json:"location,omitempty"`    // Event location

    // Optional commonly used properties
    URL        string   `json:"url,omitempty"`        // Associated URL
    Status     string   `json:"status,omitempty"`     // Event status (TENTATIVE, CONFIRMED, CANCELLED)
    Categories []string `json:"categories,omitempty"` // Event categories

    // Classification and access
    Class  string `json:"class,omitempty"`  // Access classification (PUBLIC, PRIVATE, CONFIDENTIAL)
    Transp string `json:"transp,omitempty"` // Time transparency (OPAQUE, TRANSPARENT)

    // Organizational properties
    Organizer string   `json:"organizer,omitempty"` // Event organizer
    Attendees []string `json:"attendees,omitempty"` // Event attendees

    // Scheduling properties
    Priority int `json:"priority,omitempty"` // Priority (0-9, 0=undefined)
    Sequence int `json:"sequence,omitempty"` // Revision sequence number

    // Date/Time metadata
    Created      string `json:"created,omitempty"`       // Creation date-time
    LastModified string `json:"last_modified,omitempty"` // Last modification date-time

    // Recurrence properties
    RRule        string   `json:"rrule,omitempty"`         // Recurrence rule
    RecurrenceID string   `json:"recurrence_id,omitempty"` // Recurrence identifier
    ExDates      []string `json:"exdates,omitempty"`       // Exception dates
    RDates       []string `json:"rdates,omitempty"`        // Recurrence dates

    // Other properties
    Geo       Geolocation `json:"geo,omitempty"`        // Geographic position (latitude;longitude)
    Resources []string    `json:"resources,omitempty"`  // Resources needed
    Contact   string      `json:"contact,omitempty"`    // Contact information
    RelatedTo string      `json:"related_to,omitempty"` // Related to other component
    Comment   string      `json:"comment,omitempty"`    // Comment
}
```

<a name="Geolocation"></a>
## type [Geolocation](<https://github.com:beyondcivic/icaljson/blob/main/pkg/icaljson/structs.go#L17-L20>)



```go
type Geolocation struct {
    Latitude  float64 `json:"latitude,omitempty"`
    Longitude float64 `json:"longitude,omitempty"`
}
```