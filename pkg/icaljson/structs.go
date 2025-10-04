package icaljson

// Calendar represents a VCALENDAR component according to RFC 5545
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

// Event represents a VEVENT component according to RFC 5545
type Event struct {
	// Required properties (in most contexts)
	UID     string `json:"uid,omitempty"`     // Unique identifier
	DTStamp string `json:"dtstamp,omitempty"` // Date-time stamp

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
	Geo       string   `json:"geo,omitempty"`        // Geographic position (latitude;longitude)
	Resources []string `json:"resources,omitempty"`  // Resources needed
	Contact   string   `json:"contact,omitempty"`    // Contact information
	RelatedTo string   `json:"related_to,omitempty"` // Related to other component
	Comment   string   `json:"comment,omitempty"`    // Comment
}
