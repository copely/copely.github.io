package models

// TileServiceCategory enumerates supported Copely tiling service slugs.
type TileServiceCategory string

const (
	ServicePorcelainCeramic   TileServiceCategory = "porcelain-ceramic"
	ServiceNaturalStone       TileServiceCategory = "natural-stone"
	ServiceWetRooms           TileServiceCategory = "wet-rooms"
	ServiceKitchenBacksplash  TileServiceCategory = "kitchen-backsplashes"
	ServiceOutdoorPatios      TileServiceCategory = "outdoor-patios"
)

// ProjectTimeline enumerates preferred start windows for estimate requests.
type ProjectTimeline string

const (
	TimelineThisWeek      ProjectTimeline = "this-week"
	TimelineNext2Weeks    ProjectTimeline = "next-2-weeks"
	TimelinePlanningAhead ProjectTimeline = "planning-ahead"
)

// IsValidService verifies a service category slug matches a supported Copely service.
func IsValidService(s TileServiceCategory) bool {
	switch s {
	case ServicePorcelainCeramic,
		ServiceNaturalStone,
		ServiceWetRooms,
		ServiceKitchenBacksplash,
		ServiceOutdoorPatios:
		return true
	default:
		return false
	}
}

// IsValidTimeline verifies a timeline value matches a supported option.
func IsValidTimeline(t ProjectTimeline) bool {
	switch t {
	case TimelineThisWeek, TimelineNext2Weeks, TimelinePlanningAhead:
		return true
	default:
		return false
	}
}

// EstimateRequestPayload represents an incoming lead/estimate form submission.
type EstimateRequestPayload struct {
	FirstName string              `json:"firstName,omitempty"`
	LastName  string              `json:"lastName,omitempty"`
	Phone     string              `json:"phone,omitempty"`
	Email     string              `json:"email,omitempty"`
	Address   string              `json:"address,omitempty"`
	Town      string              `json:"town,omitempty"`
	Service   TileServiceCategory `json:"service"`
	Timeline  ProjectTimeline     `json:"timeline"`
	Notes     string              `json:"notes,omitempty"`
}

// EstimateRecord stores a persisted estimate submission with server metadata.
type EstimateRecord struct {
	ID             string                 `json:"id"`
	ResponseWindow string                 `json:"responseWindow"`
	CreatedAt      string                 `json:"createdAt"`
	Payload        EstimateRequestPayload `json:"payload"`
}

// EstimateResponse is returned to the client upon successful estimate submission.
type EstimateResponse struct {
	Success        bool                `json:"success"`
	EstimateID     string              `json:"estimateId"`
	ResponseWindow string              `json:"responseWindow"`
	Service        TileServiceCategory `json:"service"`
	Timeline       ProjectTimeline     `json:"timeline"`
	CreatedAt      string              `json:"createdAt"`
}

// BookingRequestPayload represents a requested on-site consultation slot.
type BookingRequestPayload struct {
	Date          string              `json:"date"`     // YYYY-MM-DD
	TimeSlot      string              `json:"timeSlot"` // 24h HH:MM, e.g. "08:00"
	Town          string              `json:"town"`
	Service       TileServiceCategory `json:"service,omitempty"`
	Timeline      ProjectTimeline     `json:"timeline,omitempty"`
	CustomerName  string              `json:"customerName,omitempty"`
	CustomerPhone string              `json:"customerPhone,omitempty"`
	CustomerEmail string              `json:"customerEmail,omitempty"`
	Notes         string              `json:"notes,omitempty"`
}

// BookingRecord stores a confirmed on-site booking slot.
type BookingRecord struct {
	ID        string                `json:"id"`
	CreatedAt string                `json:"createdAt"`
	Payload   BookingRequestPayload `json:"payload"`
}

// BookingResponse is returned to the client upon booking confirmation.
type BookingResponse struct {
	Success   bool                `json:"success"`
	BookingID string              `json:"bookingId"`
	Date      string              `json:"date"`
	TimeSlot  string              `json:"timeSlot"`
	Town      string              `json:"town"`
	Service   TileServiceCategory `json:"service,omitempty"`
	Timeline  ProjectTimeline     `json:"timeline,omitempty"`
	CreatedAt string              `json:"createdAt"`
}

// TimeSlotAvailability models a single 30-minute consultation slot.
type TimeSlotAvailability struct {
	Time      string `json:"time"`      // "08:00"
	Available bool   `json:"available"` // false if already booked
}

// DayAvailability describes slot availability for one calendar date.
type DayAvailability struct {
	Date      string                 `json:"date"`      // YYYY-MM-DD
	DayOfWeek int                    `json:"dayOfWeek"` // 0 (Sun) - 6 (Sat)
	Available bool                   `json:"available"`
	Slots     []TimeSlotAvailability `json:"slots"`
}

// BookingAvailabilityResponse is returned for GET /api/bookings/availability.
type BookingAvailabilityResponse struct {
	Month string            `json:"month"` // YYYY-MM
	Days  []DayAvailability `json:"days"`
}

// ServiceItem models a Copely service offering for cards and detail pages.
type ServiceItem struct {
	ID                string              `json:"id"`
	Slug              TileServiceCategory `json:"slug"`
	Title             string              `json:"title"`
	ShortSubtitle     string              `json:"shortSubtitle"`
	Description       string              `json:"description"`
	ImageURL          string              `json:"imageUrl"`
	WarrantyYears     int                 `json:"warrantyYears"`
	TimelineLabel     string              `json:"timelineLabel"`
	IdealFor          string              `json:"idealFor"`
	Specs             []string            `json:"specs"`
	ProcessHighlights []string            `json:"processHighlights"`
}

// BeforeAfterCase models a matched before/after restoration pair.
type BeforeAfterCase struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Location    string `json:"location"`
	BeforeImage string `json:"beforeImage"`
	AfterImage  string `json:"afterImage"`
	Featured    bool   `json:"featured"`
}

// TestimonialMetric represents one summary metric on a customer case study card.
type TestimonialMetric struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Testimonial models a featured customer story.
type Testimonial struct {
	ID           string              `json:"id"`
	Author       string              `json:"author"`
	ProjectType  string              `json:"projectType"`
	Location     string              `json:"location"`
	AvatarURL    string              `json:"avatarUrl"`
	Quote        string              `json:"quote"`
	HeroImageURL string              `json:"heroImageUrl"`
	Metrics      []TestimonialMetric `json:"metrics"`
}

// RecentWorkItem models an entry in the Recent Work masonry/bento grid.
type RecentWorkItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Location string `json:"location"`
	Category string `json:"category"`
	ImageURL string `json:"imageUrl"`
}

// FaqItem represents a single question and answer entry.
type FaqItem struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// FaqCategoryGroup organizes FAQ items by topic for the /faq page.
type FaqCategoryGroup struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Items []FaqItem `json:"items"`
}

// GoogleReview represents a verified customer review card on /faq.
type GoogleReview struct {
	ID         string `json:"id"`
	Author     string `json:"author"`
	Initials   string `json:"initials"`
	Rating     int    `json:"rating"`
	DateLabel  string `json:"dateLabel"`
	ProjectTag string `json:"projectTag"`
	Body       string `json:"body"`
}

// GlossaryTerm defines a tiling term on the /faq page.
type GlossaryTerm struct {
	ID         string `json:"id"`
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

// BlogSection models one structured section inside a full blog article.
type BlogSection struct {
	Heading    string   `json:"heading"`
	Paragraphs []string `json:"paragraphs"`
}

// BlogPost represents a Copely editorial article for /blog and /blog/:slug.
type BlogPost struct {
	ID           string        `json:"id"`
	Slug         string        `json:"slug"`
	Category     string        `json:"category"`
	PublishedAt  string        `json:"publishedAt"`
	ReadTime     string        `json:"readTime"`
	Title        string        `json:"title"`
	Excerpt      string        `json:"excerpt"`
	AuthorName   string        `json:"authorName"`
	AuthorRole   string        `json:"authorRole"`
	AuthorAvatar string        `json:"authorAvatar"`
	CoverImage   string        `json:"coverImage"`
	Featured     bool          `json:"featured"`
	Sections     []BlogSection `json:"sections"`
}

// ContentBundle aggregates all seeded public content returned by GET /api/content.
type ContentBundle struct {
	Services         []ServiceItem      `json:"services"`
	BeforeAfterCases []BeforeAfterCase  `json:"beforeAfterCases"`
	Testimonials     []Testimonial      `json:"testimonials"`
	RecentWork       []RecentWorkItem   `json:"recentWork"`
	HomeFaqs         []FaqItem          `json:"homeFaqs"`
	FaqCategories    []FaqCategoryGroup `json:"faqCategories"`
	GoogleReviews    []GoogleReview     `json:"googleReviews"`
	GlossaryTerms    []GlossaryTerm     `json:"glossaryTerms"`
	BlogPosts        []BlogPost         `json:"blogPosts"`
	Towns            []string           `json:"towns"`
}
