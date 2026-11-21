package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"copely-backend/internal/models"
)

// ErrSlotAlreadyBooked is returned when a booking slot is no longer available.
var ErrSlotAlreadyBooked = errors.New("requested date and time slot is already booked")

// ErrInvalidBookingDate is returned when a booking date cannot be booked.
var ErrInvalidBookingDate = errors.New("requested date is unavailable for on-site visits")

// ErrNotFound is returned when a requested resource slug does not exist.
var ErrNotFound = errors.New("resource not found")

var standardSlots = []string{
	"08:00", "08:30", "09:00", "09:30",
	"10:00", "10:30", "11:00", "11:30",
	"13:00", "13:30", "14:00", "14:30",
	"15:00", "15:30", "16:00", "16:30",
}

type submissionSnapshot struct {
	Estimates []models.EstimateRecord `json:"estimates"`
	Bookings  []models.BookingRecord  `json:"bookings"`
}

// MemoryStore is a thread-safe in-memory store backed by atomic JSON persistence.
type MemoryStore struct {
	mu           sync.RWMutex
	snapshotPath string
	estimates    []models.EstimateRecord
	bookings     []models.BookingRecord
	bookedSlots  map[string]bool
	content      models.ContentBundle
	seq          int
	nowFunc      func() time.Time
}

// NewMemoryStore initializes seeded Copely content and loads any persisted snapshot.
func NewMemoryStore(snapshotPath string) (*MemoryStore, error) {
	s := &MemoryStore{
		snapshotPath: snapshotPath,
		bookedSlots:  make(map[string]bool),
		content:      seedCopelyContent(),
		nowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Pre-populate realistic already-reserved consultation slots in November 2026
	// while keeping 2026-11-18 morning slots open per the specification.
	seedBooked := []string{
		"2026-11-18|11:00",
		"2026-11-18|14:30",
		"2026-11-19|08:30",
		"2026-11-19|10:00",
		"2026-11-20|09:00",
		"2026-11-23|13:30",
		"2026-11-24|15:00",
	}
	for _, key := range seedBooked {
		s.bookedSlots[key] = true
	}

	if snapshotPath != "" {
		if err := s.loadSnapshot(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *MemoryStore) loadSnapshot() error {
	data, err := os.ReadFile(s.snapshotPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read submissions snapshot: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil
	}

	var snap submissionSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return fmt.Errorf("decode submissions snapshot: %w", err)
	}

	s.estimates = snap.Estimates
	s.bookings = snap.Bookings
	for _, b := range snap.Bookings {
		s.bookedSlots[slotKey(b.Payload.Date, b.Payload.TimeSlot)] = true
	}
	s.seq = len(s.estimates) + len(s.bookings)
	return nil
}

// persistLocked writes estimates and bookings atomically to snapshotPath.
// Caller must hold s.mu write lock.
func (s *MemoryStore) persistLocked() error {
	if s.snapshotPath == "" {
		return nil
	}

	dir := filepath.Dir(s.snapshotPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create snapshot directory: %w", err)
	}

	snap := submissionSnapshot{
		Estimates: s.estimates,
		Bookings:  s.bookings,
	}
	payload, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal submissions snapshot: %w", err)
	}

	tmp, err := os.CreateTemp(dir, "submissions-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp snapshot: %w", err)
	}
	tmpName := tmp.Name()

	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}

	if _, err := tmp.Write(payload); err != nil {
		cleanup()
		return fmt.Errorf("write temp snapshot: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		cleanup()
		return fmt.Errorf("sync temp snapshot: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp snapshot: %w", err)
	}
	if err := os.Rename(tmpName, s.snapshotPath); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename temp snapshot: %w", err)
	}
	return nil
}

func slotKey(date, timeSlot string) string {
	return date + "|" + timeSlot
}

func isValidSlot(slot string) bool {
	for _, candidate := range standardSlots {
		if candidate == slot {
			return true
		}
	}
	return false
}

// SaveEstimate stores a validated estimate submission and writes the snapshot.
func (s *MemoryStore) SaveEstimate(payload models.EstimateRequestPayload) (models.EstimateRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	now := s.nowFunc()
	rec := models.EstimateRecord{
		ID:             fmt.Sprintf("EST-%s-%04d", now.Format("2006"), 1000+s.seq),
		ResponseWindow: "Within 24 hrs",
		CreatedAt:      now.Format(time.RFC3339),
		Payload:        payload,
	}

	s.estimates = append(s.estimates, rec)
	if err := s.persistLocked(); err != nil {
		s.estimates = s.estimates[:len(s.estimates)-1]
		s.seq--
		return models.EstimateRecord{}, err
	}
	return rec, nil
}

// GetBookingAvailability computes date and 30-minute slot availability for YYYY-MM.
func (s *MemoryStore) GetBookingAvailability(month string) (models.BookingAvailabilityResponse, error) {
	parsedMonth, err := time.Parse("2006-01", month)
	if err != nil {
		return models.BookingAvailabilityResponse{}, fmt.Errorf("invalid month format %q (expected YYYY-MM)", month)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Earliest bookable date anchors to Nov 18, 2026 (matching local demo environment)
	// or today's UTC date if later.
	anchor := time.Date(2026, 11, 18, 0, 0, 0, 0, time.UTC)
	today := s.nowFunc().Truncate(24 * time.Hour)
	if today.Before(anchor) {
		// Keep Nov 18, 2026 as minimum active day when system clock matches anchor.
		anchor = today
	}

	firstOfMonth := time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := firstOfMonth.AddDate(0, 1, 0)

	var days []models.DayAvailability
	for d := firstOfMonth; d.Before(nextMonth); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		weekday := int(d.Weekday())

		// Sundays are closed; dates prior to anchor are unavailable.
		isOpenDay := weekday != int(time.Sunday) && !d.Before(anchor)
		// Thanksgiving Day 2026-11-26 is closed.
		if dateStr == "2026-11-26" {
			isOpenDay = false
		}

		slots := make([]models.TimeSlotAvailability, 0, len(standardSlots))
		openCount := 0
		for _, t := range standardSlots {
			available := isOpenDay && !s.bookedSlots[slotKey(dateStr, t)]
			if available {
				openCount++
			}
			slots = append(slots, models.TimeSlotAvailability{
				Time:      t,
				Available: available,
			})
		}

		days = append(days, models.DayAvailability{
			Date:      dateStr,
			DayOfWeek: weekday,
			Available: isOpenDay && openCount > 0,
			Slots:     slots,
		})
	}

	return models.BookingAvailabilityResponse{
		Month: month,
		Days:  days,
	}, nil
}

// BookSlot atomically reserves an available slot and persists the snapshot.
func (s *MemoryStore) BookSlot(payload models.BookingRequestPayload) (models.BookingRecord, error) {
	parsedDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return models.BookingRecord{}, ErrInvalidBookingDate
	}
	if parsedDate.Weekday() == time.Sunday || payload.Date == "2026-11-26" {
		return models.BookingRecord{}, ErrInvalidBookingDate
	}
	if !isValidSlot(payload.TimeSlot) {
		return models.BookingRecord{}, fmt.Errorf("invalid time slot %q", payload.TimeSlot)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := slotKey(payload.Date, payload.TimeSlot)
	if s.bookedSlots[key] {
		return models.BookingRecord{}, ErrSlotAlreadyBooked
	}

	s.seq++
	now := s.nowFunc()
	rec := models.BookingRecord{
		ID:        fmt.Sprintf("BK-%s-%04d", now.Format("2006"), 1000+s.seq),
		CreatedAt: now.Format(time.RFC3339),
		Payload:   payload,
	}

	s.bookedSlots[key] = true
	s.bookings = append(s.bookings, rec)
	if err := s.persistLocked(); err != nil {
		delete(s.bookedSlots, key)
		s.bookings = s.bookings[:len(s.bookings)-1]
		s.seq--
		return models.BookingRecord{}, err
	}
	return rec, nil
}

// GetContent returns a copy of the seeded ContentBundle.
func (s *MemoryStore) GetContent() models.ContentBundle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.content
}

// GetServiceBySlug returns a single service matching slug.
func (s *MemoryStore) GetServiceBySlug(slug string) (models.ServiceItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, svc := range s.content.Services {
		if string(svc.Slug) == slug {
			return svc, nil
		}
	}
	return models.ServiceItem{}, ErrNotFound
}

// GetBlogPostBySlug returns a single blog post matching slug.
func (s *MemoryStore) GetBlogPostBySlug(slug string) (models.BlogPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, post := range s.content.BlogPosts {
		if post.Slug == slug {
			return post, nil
		}
	}
	return models.BlogPost{}, ErrNotFound
}

func seedCopelyContent() models.ContentBundle {
	services := []models.ServiceItem{
		{
			ID:            "svc-porcelain",
			Slug:          models.ServicePorcelainCeramic,
			Title:         "Porcelain & Tile",
			ShortSubtitle: "Large-format floors & walls",
			Description:   "Laser-leveled subfloor prep, uncoupling membranes, and large-format porcelain or ceramic tile set with razor-thin lippage-free grout joints.",
			ImageURL:      "/images/service-porcelain.jpg",
			WarrantyYears: 25,
			TimelineLabel: "2–4 days typical",
			IdealFor:      "Open-concept living rooms, kitchens, foyers, and high-traffic hallways.",
			Specs: []string{
				"ANSI A108 substrate flattening to < 1/8\" over 10 ft",
				"Crack-isolation uncoupling membrane over concrete slabs",
				"100% coverage back-buttering on tiles larger than 15\"",
				"Stain-proof high-performance epoxy or urethane grout",
			},
			ProcessHighlights: []string{
				"Dust-extracted demo & slab grinding",
				"Self-leveling underlayment & movement joints",
				"Laser layout for balanced perimeter cuts",
			},
		},
		{
			ID:            "svc-wet-rooms",
			Slug:          models.ServiceWetRooms,
			Title:         "Wet Rooms",
			ShortSubtitle: "Seamless curbless showers & waterproofing",
			Description:   "Flood-tested bonded waterproof membranes, pre-sloped shower pans, linear drains, and custom tiled niches engineered to never leak.",
			ImageURL:      "/images/service-wetroom.jpg",
			WarrantyYears: 25,
			TimelineLabel: "3–5 days typical",
			IdealFor:      "Primary suite walk-in showers, curbless accessible wet rooms, and steam enclosures.",
			Specs: []string{
				"Continuous sheet-membrane waterproofing wall-to-floor",
				"24-hour standing water flood test documented on-site",
				"Precision 1/4\" per foot pan slope to linear or square drains",
				"100% silicone movement joints at all change-of-plane corners",
			},
			ProcessHighlights: []string{
				"Framing inspection & blocking for benches/glass",
				"Bonded pipe boots & mixing valve collar seals",
				"Mitered stone edges & recessed niche waterproofing",
			},
		},
		{
			ID:            "svc-backsplashes",
			Slug:          models.ServiceKitchenBacksplash,
			Title:         "Backsplashes",
			ShortSubtitle: "Custom kitchen mosaics, zellige & stone",
			Description:   "Artisan zellige, herringbone, and slab backsplashes finished around outlets, range hoods, and custom millwork with exact symmetry.",
			ImageURL:      "/images/service-backsplash.jpg",
			WarrantyYears: 25,
			TimelineLabel: "1–2 days typical",
			IdealFor:      "Chef kitchens, butler's pantries, fireplace surrounds, and wet bars.",
			Specs: []string{
				"Flush outlet box extensions & laser-centered focal layouts",
				"Color-matched anodized metal or mitered edge returns",
				"Non-sag polymer-modified mortar for handmade clay & glass",
				"Penetrating sealer applied prior to and after grouting",
			},
			ProcessHighlights: []string{
				"Full countertop & range protection before cut #1",
				"Dry-lay pattern blending for tonal variation",
				"Clean silicone perimeter caulk along countertops",
			},
		},
		{
			ID:            "svc-natural-stone",
			Slug:          models.ServiceNaturalStone,
			Title:         "Natural Stone",
			ShortSubtitle: "Marble, travertine, slate & limestone",
			Description:   "Honed limestone, veined Carrara marble, and Hill Country travertine installed with non-staining white thin-set and penetrating sealers.",
			ImageURL:      "/images/service-natural-stone.jpg",
			WarrantyYears: 25,
			TimelineLabel: "3–5 days typical",
			IdealFor:      "Architectural foyers, dining rooms, fireplace hearths, and primary baths.",
			Specs: []string{
				"Crack-suppressing membrane beneath moisture-sensitive stone",
				"Non-staining white rapid-setting stone mortar",
				"Vein-matched dry layout approved with homeowner",
				"Breathable fluoropolymer penetrating stone sealer",
			},
			ProcessHighlights: []string{
				"Zero-deflection subfloor stiffening check",
				"Wet diamond polishing on exposed stone edges",
				"Care guide & neutral-pH maintenance kit included",
			},
		},
		{
			ID:            "svc-outdoor-patios",
			Slug:          models.ServiceOutdoorPatios,
			Title:         "Outdoor Patios",
			ShortSubtitle: "Frost-proof pavers, pool decks & steps",
			Description:   "Exterior porcelain pavers, natural stone terraces, and pool copings engineered with positive drainage and freeze-thaw resistant mortars.",
			ImageURL:      "/images/service-outdoor-patio.jpg",
			WarrantyYears: 25,
			TimelineLabel: "3–6 days typical",
			IdealFor:      "Covered porches, pool surrounds, outdoor kitchens, and hillside stone steps.",
			Specs: []string{
				"Waterproofing & anti-efflorescence drainage mat",
				"Freeze-thaw rated exterior polymer mortar & epoxy grout",
				"Engineered perimeter expansion joints for Texas heat swings",
				"Slip-resistant R11 textured stone & exterior porcelain finishes",
			},
			ProcessHighlights: []string{
				"Slope verification away from home foundation",
				"Custom bullnose & mitered step treads",
				"Hydrophobic surface protection against oak tannins",
			},
		},
	}

	beforeAfterCases := []models.BeforeAfterCase{
		{
			ID:          "ba-hero",
			Title:       "Open-Concept Porcelain Slab Replacement",
			Subtitle:    "Cracked 1990s builder ceramic replaced with 24×48 charcoal slate porcelain",
			Location:    "Dripping Springs, TX",
			BeforeImage: "/images/ba-hero-before.jpg",
			AfterImage:  "/images/ba-hero-after.jpg",
			Featured:    true,
		},
		{
			ID:          "ba-bath",
			Title:       "Primary Marble Wet Room",
			Subtitle:    "Leaking fiberglass-pan shower rebuilt into curbless Carrara marble suite",
			Location:    "Bee Cave, TX",
			BeforeImage: "/images/ba-bath-before.jpg",
			AfterImage:  "/images/ba-bath-after.jpg",
			Featured:    false,
		},
		{
			ID:          "ba-patio",
			Title:       "Limestone & Slate Paver Terrace",
			Subtitle:    "Weathered bare concrete deck transformed into custom Hill Country stone patio",
			Location:    "Wimberley, TX",
			BeforeImage: "/images/ba-patio-before.jpg",
			AfterImage:  "/images/ba-patio-after.jpg",
			Featured:    false,
		},
		{
			ID:          "ba-kitchen",
			Title:       "Artisan Glazed Herringbone Splashback",
			Subtitle:    "Peeling dated backsplash replaced with handmade glazed herringbone tile",
			Location:    "Lakeway, TX",
			BeforeImage: "/images/ba-kitchen-before.jpg",
			AfterImage:  "/images/ba-kitchen-after.jpg",
			Featured:    false,
		},
	}

	testimonials := []models.Testimonial{
		{
			ID:           "test-daniel",
			Author:       "Daniel Okafor",
			ProjectType:  "Tile repair",
			Location:     "Houston, TX",
			AvatarURL:    "/images/avatar-1.jpg",
			Quote:        "Two other companies told me I needed a full slab replacement. Copely found one bad expansion joint, fixed it in an afternoon, and charged me a few hundred dollars. They talked themselves out of a much bigger job.",
			HeroImageURL: "/images/case-1.jpg",
			Metrics: []models.TestimonialMetric{
				{Label: "Turnaround", Value: "1 Afternoon"},
				{Label: "Transparency", Value: "2 Quotes beaten"},
				{Label: "Protection", Value: "25yr Warranty"},
			},
		},
		{
			ID:           "test-claire",
			Author:       "Claire Morales",
			ProjectType:  "Whole-home porcelain",
			Location:     "Dripping Springs, TX",
			AvatarURL:    "/images/avatar-2.jpg",
			Quote:        "They leveled our slab the same week we called, and 1,200 sq ft of new large-format porcelain was down four days later. Every transition at the patio doors is dead flat.",
			HeroImageURL: "/images/case-2.jpg",
			Metrics: []models.TestimonialMetric{
				{Label: "Scope", Value: "1,200 sq ft"},
				{Label: "Schedule", Value: "4 Days"},
				{Label: "Protection", Value: "25yr Warranty"},
			},
		},
		{
			ID:           "test-harold",
			Author:       "Harold & Beth Vance",
			ProjectType:  "Curbless wet room",
			Location:     "Fredericksburg, TX",
			AvatarURL:    "/images/avatar-3.jpg",
			Quote:        "After a slow leak ruined our downstairs ceiling, Copely stripped the shower to studs, flood-tested the new pan for 24 hours, and set bookmatched stone in five days flat.",
			HeroImageURL: "/images/case-3.jpg",
			Metrics: []models.TestimonialMetric{
				{Label: "Flood Test", Value: "24hr Verified"},
				{Label: "Duration", Value: "5 Days"},
				{Label: "Protection", Value: "25yr Warranty"},
			},
		},
		{
			ID:           "test-nadia",
			Author:       "Nadia Al-Mansoor",
			ProjectType:  "Zellige kitchen backsplash",
			Location:     "Austin, TX",
			AvatarURL:    "/images/avatar-4.jpg",
			Quote:        "Our handmade Moroccan tile was tricky because every edge varies slightly. Sam's crew dry-laid the entire wall before touching mortar, and the finish around our brass hood is museum quality.",
			HeroImageURL: "/images/case-4.jpg",
			Metrics: []models.TestimonialMetric{
				{Label: "Duration", Value: "1.5 Days"},
				{Label: "Edge Finish", Value: "Mitered Trim"},
				{Label: "Protection", Value: "25yr Warranty"},
			},
		},
		{
			ID:           "test-graham",
			Author:       "Graham Whitaker",
			ProjectType:  "Outdoor limestone patio",
			Location:     "Wimberley, TX",
			AvatarURL:    "/images/avatar-5.jpg",
			Quote:        "Texas summers and heavy storms wrecked our old flagstone mortar. Copely installed an uncoupling drainage mat and frost-proof stone pavers that still look brand new two winters later.",
			HeroImageURL: "/images/case-5.jpg",
			Metrics: []models.TestimonialMetric{
				{Label: "Area", Value: "680 sq ft"},
				{Label: "System", Value: "Drainage Mat"},
				{Label: "Protection", Value: "25yr Warranty"},
			},
		},
	}

	recentWork := []models.RecentWorkItem{
		{
			ID:       "rw-1",
			Title:    "Notched Mortar Bed & Membrane Prep",
			Location: "Austin, TX",
			Category: "Substrate Prep",
			ImageURL: "/images/bento-1.jpg",
		},
		{
			ID:       "rw-2",
			Title:    "White Oak Herringbone Porcelain",
			Location: "Dripping Springs, TX",
			Category: "Porcelain & Tile",
			ImageURL: "/images/bento-2.jpg",
		},
		{
			ID:       "rw-3",
			Title:    "Custom Stone Bench & Niche Wet Room",
			Location: "Bee Cave, TX",
			Category: "Wet Rooms",
			ImageURL: "/images/bento-3.jpg",
		},
		{
			ID:       "rw-4",
			Title:    "Limestone Pool Coping & Terrace Steps",
			Location: "Wimberley, TX",
			Category: "Outdoor Patios",
			ImageURL: "/images/bento-4.jpg",
		},
		{
			ID:       "rw-5",
			Title:    "Precision Lippage-Free Leveling System",
			Location: "Lakeway, TX",
			Category: "Large-Format",
			ImageURL: "/images/bento-5.jpg",
		},
		{
			ID:       "rw-6",
			Title:    "Bookmatched Gallery Entryway Walls",
			Location: "Fredericksburg, TX",
			Category: "Natural Stone",
			ImageURL: "/images/bento-6.jpg",
		},
	}

	homeFaqs := []models.FaqItem{
		{
			ID:       "faq-1",
			Category: "estimates",
			Question: "How much does an estimate cost?",
			Answer:   "Nothing. We measure your space, take photos, check your subfloor or shower pan, and give you a clear written estimate — whether or not you hire us.",
		},
		{
			ID:       "faq-2",
			Category: "during-job",
			Question: "How long does a full tile installation take?",
			Answer:   "Kitchen backsplashes typically take 1–2 days, most floor installations finish in 2–4 days, and complete waterproofed curbless wet rooms take 3–5 days including our 24-hour pan flood test.",
		},
		{
			ID:       "faq-3",
			Category: "estimates",
			Question: "Do you handle water damage or insurance claims?",
			Answer:   "Yes. When a leaking shower pan or broken pipe damages subfloors or tile, we document moisture readings, supply line-item estimates compatible with major carriers, and bill insurance directly where approved.",
		},
		{
			ID:       "faq-4",
			Category: "warranty",
			Question: "What does the warranty actually cover?",
			Answer:   "Our 25-year transferable warranty covers waterproofing integrity, mortar bond failure, substrate prep, lippage tolerances, and structural grout cracking — backed in writing for you and the next homeowner.",
		},
		{
			ID:       "faq-5",
			Category: "waterproofing",
			Question: "Can you waterproof my shower before tiling?",
			Answer:   "Always. Tile and grout themselves are never waterproof. We install continuous bonded sheet waterproofing membranes, sealed valve collars, and sloped mortar beds before setting a single tile.",
		},
		{
			ID:       "faq-6",
			Category: "during-job",
			Question: "What happens to all the old tile debris and dust?",
			Answer:   "We mask walkways with heavy floor protection, run HEPA dust-extraction shrouds on our grinders and demo hammers, and haul every scrap of old tile and mortar off-site daily.",
		},
	}

	faqCategories := []models.FaqCategoryGroup{
		{
			ID:    "estimates",
			Title: "Estimates & pricing",
			Items: []models.FaqItem{
				homeFaqs[0],
				{
					ID:       "faq-est-2",
					Category: "estimates",
					Question: "Can you supply the tile or work with tile I already bought?",
					Answer:   "Both. You can purchase tile through our Hill Country trade showrooms at contractor pricing, or supply your own material and we bring all membranes, mortar, trims, and grout.",
				},
				{
					ID:       "faq-est-3",
					Category: "estimates",
					Question: "Are there hidden fees if my slab isn't level?",
					Answer:   "No surprises. We laser-check floor flatness during your free estimate visit so self-leveling compound or joist stiffening is priced transparently upfront.",
				},
				homeFaqs[2],
			},
		},
		{
			ID:    "during-job",
			Title: "During the job",
			Items: []models.FaqItem{
				homeFaqs[1],
				homeFaqs[5],
				{
					ID:       "faq-job-3",
					Category: "during-job",
					Question: "Can we stay in our home while you tile?",
					Answer:   "Yes. We set up zippered dust containment barriers, clean up tools at the end of every shift, and plan room sequencing so your kitchen or second bath stays usable.",
				},
				{
					ID:       "faq-job-4",
					Category: "during-job",
					Question: "When can we walk on new floors or turn on the shower?",
					Answer:   "Light foot traffic is usually safe 12–18 hours after grouting. Wet rooms and steam showers are ready for full use 48 hours after final silicone curing.",
				},
			},
		},
		{
			ID:    "system",
			Title: "Your tile system, explained",
			Items: []models.FaqItem{
				{
					ID:       "faq-sys-1",
					Category: "system",
					Question: "Why do you use an uncoupling membrane over concrete?",
					Answer:   "Central Texas clay soil expands and contracts seasonally, creating hairline slab movement. An uncoupling membrane absorbs lateral shear so your porcelain or stone never cracks.",
				},
				{
					ID:       "faq-sys-2",
					Category: "system",
					Question: "What makes large-format tile harder to install?",
					Answer:   "Tiles 15 inches or longer amplify minor subfloor dips and require 95%+ mortar coverage via back-buttering plus mechanical leveling clips to guarantee zero lippage.",
				},
				{
					ID:       "faq-sys-3",
					Category: "system",
					Question: "Do you use cement grout or epoxy grout?",
					Answer:   "We recommend stain-proof high-performance urethane or epoxy grout in showers and kitchens so coffee, hard water, and soap never discolor grout joints.",
				},
			},
		},
		{
			ID:    "waterproofing",
			Title: "Waterproofing & prep",
			Items: []models.FaqItem{
				homeFaqs[4],
				{
					ID:       "faq-wp-2",
					Category: "waterproofing",
					Question: "Do you flood-test every new shower pan?",
					Answer:   "Yes. Every custom wet room pan is plugged and filled for 24 hours to verify zero water drop before a single floor tile goes down.",
				},
				{
					ID:       "faq-wp-3",
					Category: "waterproofing",
					Question: "How do you prevent mold in shower corners?",
					Answer:   "Instead of rigid grout at inside wall-to-floor corners—which always cracks—we bridge transitions with waterproof banding and finish with 100% mold-resistant silicone.",
				},
			},
		},
		{
			ID:    "warranty",
			Title: "Warranty & after",
			Items: []models.FaqItem{
				homeFaqs[3],
				{
					ID:       "faq-war-2",
					Category: "warranty",
					Question: "Is the 25-year warranty transferable if I sell my house?",
					Answer:   "Yes. Your signed warranty certificate transfers automatically to the next owner, adding verifiable value during home inspections.",
				},
				{
					ID:       "faq-war-3",
					Category: "warranty",
					Question: "Do natural stone floors need periodic resealing?",
					Answer:   "Interior honed marble and limestone typically need a quick wipe-on penetrating sealer every 3–5 years; exterior stone patios benefit from sealing every 2–3 years.",
				},
			},
		},
	}

	googleReviews := []models.GoogleReview{
		{
			ID:         "rev-1",
			Author:     "Max Rivera",
			Initials:   "MR",
			Rating:     5,
			DateLabel:  "2 weeks ago",
			ProjectTag: "Porcelain Living Floor",
			Body:       "Copely ground down a high concrete hump in our living room before setting 24x48 porcelain. Zero hollow spots, zero lippage, and their crew swept the driveway before leaving each day.",
		},
		{
			ID:         "rev-2",
			Author:     "J. Whitfield",
			Initials:   "JW",
			Rating:     5,
			DateLabel:  "1 month ago",
			ProjectTag: "Curbless Wet Room",
			Body:       "Night-and-day difference from our old contractor. They caught missing stud blocking, flood-tested the shower pan overnight, and showed us photos of every waterproofing layer before tiling.",
		},
		{
			ID:         "rev-3",
			Author:     "Gary Hutchins",
			Initials:   "GH",
			Rating:     5,
			DateLabel:  "2 months ago",
			ProjectTag: "Outdoor Pool Patio",
			Body:       "Jordan walked us through drainage mats and exterior frost-proof pavers for our Wimberley deck. Finished right on the quoted price even after fixing a tricky threshold step.",
		},
		{
			ID:         "rev-4",
			Author:     "Janet Solomon",
			Initials:   "JS",
			Rating:     5,
			DateLabel:  "3 months ago",
			ProjectTag: "Kitchen Zellige Backsplash",
			Body:       "We worried our handmade zellige tiles would look wavy in a bad way. Sam's setters blended boxes carefully and mitered every outside corner. Looks like an architectural magazine.",
		},
		{
			ID:         "rev-5",
			Author:     "Trevor & Maya Lin",
			Initials:   "TL",
			Rating:     5,
			DateLabel:  "4 months ago",
			ProjectTag: "Whole-Home Stone",
			Body:       "Three estimates in, Copely was the only crew who checked slab moisture and expansion joints with real instruments. Best home investment we have made.",
		},
		{
			ID:         "rev-6",
			Author:     "Elena Vance-Ramos",
			Initials:   "EV",
			Rating:     5,
			DateLabel:  "5 months ago",
			ProjectTag: "Shower Pan Repair",
			Body:       "Honest, fast, and tidy. They saved our primary bath from a subfloor rot disaster and kept dust completely out of our bedroom closet.",
		},
	}

	glossaryTerms := []models.GlossaryTerm{
		{
			ID:         "gl-1",
			Term:       "Subfloor",
			Definition: "The structural concrete slab or plywood decking beneath your tile that must be flat, rigid, and deflection-free before any setting begins.",
		},
		{
			ID:         "gl-2",
			Term:       "Uncoupling membrane",
			Definition: "A flexible engineered mat placed between the subfloor and tile so tiny seasonal slab shifts never telegraph up as cracked tiles.",
		},
		{
			ID:         "gl-3",
			Term:       "Flashing / Profiles",
			Definition: "Precision extruded aluminum or stainless edge trims that lock down exposed tile borders and protect corners from chipping.",
		},
		{
			ID:         "gl-4",
			Term:       "Drip edge / Threshold",
			Definition: "A pitched stone sill or transition saddle engineered to shed water back into the shower and bridge differing floor heights cleanly.",
		},
		{
			ID:         "gl-5",
			Term:       "Shower pan slope",
			Definition: "The exact 1/4-inch-per-foot pitch built into the mortar bed so every drop of shower water flows reliably toward the drain.",
		},
		{
			ID:         "gl-6",
			Term:       "Pipe boot / Mixing valve seal",
			Definition: "Prefabricated waterproof collars sealed directly over plumbing penetrations so steam and pressurized spray stay behind the tile.",
		},
		{
			ID:         "gl-7",
			Term:       "Expansion joint",
			Definition: "A flexible perimeter or field seam filled with matching silicone rather than rigid grout so tile expanses breathe safely with temperature swings.",
		},
		{
			ID:         "gl-8",
			Term:       "Back-buttering",
			Definition: "Skimming a flat coat of thin-set mortar onto the back of every large tile before setting it into the combed floor bed for 100% bond coverage.",
		},
		{
			ID:         "gl-9",
			Term:       "Epoxy grout",
			Definition: "A non-porous, chemical-resistant two-part grout that never needs sealing and resists hard water, wine, coffee, and mildew stains permanently.",
		},
	}

	blogPosts := []models.BlogPost{
		{
			ID:           "post-1",
			Slug:         "tile-installation-cost-texas-hill-country",
			Category:     "PRICING",
			PublishedAt:  "JUN 18, 2026",
			ReadTime:     "6 MIN",
			Title:        "What does new tile installation cost in the Hill Country?",
			Excerpt:      "A line-by-line look at subfloor flattening, waterproofing membranes, mortar specs, and labor rates across Austin, Dripping Springs, and Fredericksburg.",
			AuthorName:   "Jordan Pike",
			AuthorRole:   "Lead Estimator",
			AuthorAvatar: "/images/author-jordan.jpg",
			CoverImage:   "/images/blog-featured.jpg",
			Featured:     true,
			Sections: []models.BlogSection{
				{
					Heading: "Why two tile quotes rarely price the same system",
					Paragraphs: []string{
						"When homeowners collect three estimates for a primary bath or 1,000-square-foot living room floor, quotes can swing by 40% before anyone picks out a tile. Almost every gap hides beneath the surface: whether the installer includes slab grinding, crack-isolation membranes, and flood-tested waterproofing.",
						"At Copely, we itemize subfloor prep, uncoupling membranes, trims, and grout separately from tile selection so you know exactly what protects your floor for the next 25 years.",
					},
				},
				{
					Heading: "Typical 2026 installed ranges by project type",
					Paragraphs: []string{
						"Kitchen backsplashes in Austin and Dripping Springs generally fall between $2,200 and $4,800 depending on layout geometry (straight stack vs. herringbone) and outlet cuts. Large-format porcelain floor installations typically range from $14 to $24 per square foot installed including uncoupling membrane.",
						"Complete curbless wet room showers—including demolition, framing corrections, bonded sheet waterproofing, 24-hour flood testing, and custom stone benches—typically land between $7,500 and $14,000.",
					},
				},
				{
					Heading: "Three questions to verify before signing any bid",
					Paragraphs: []string{
						"Always ask whether your installer checks concrete flatness with a 10-foot straightedge, installs soft silicone movement joints at walls instead of hard grout, and backs waterproofing in writing.",
					},
				},
			},
		},
		{
			ID:           "post-2",
			Slug:         "5-signs-shower-pan-hidden-water-damage",
			Category:     "WATERPROOFING",
			PublishedAt:  "MAY 27, 2026",
			ReadTime:     "5 MIN",
			Title:        "5 signs your shower pan has hidden water damage",
			Excerpt:      "Musty baseboard odors, hollow tiles near the drain, and white efflorescence along curb joints are early warnings before subfloor rot starts.",
			AuthorName:   "Sam Calloway",
			AuthorRole:   "Founder & Master Setter",
			AuthorAvatar: "/images/author-sam.jpg",
			CoverImage:   "/images/blog-waterproofing.jpg",
			Featured:     false,
			Sections: []models.BlogSection{
				{
					Heading: "Water travels quietly long before a ceiling stain appears",
					Paragraphs: []string{
						"Most failed showers do not burst overnight. Capillary moisture seeps past cracked corner grout or unsealed curb fasteners into framing months before drywall shows discoloration.",
					},
				},
				{
					Heading: "What to inspect around your shower today",
					Paragraphs: []string{
						"Look for swollen baseboards outside the glass door, persistent damp grout lines 24 hours after showering, chalky mineral crust (efflorescence) along the curb, hollow-sounding floor tiles when tapped, or soft drywall beside the valve trim.",
						"Catching pan failure early often limits repairs to the wet enclosure itself rather than adjacent hardwood or framing.",
					},
				},
			},
		},
		{
			ID:           "post-3",
			Slug:         "porcelain-vs-natural-stone-texas-homes",
			Category:     "MATERIALS",
			PublishedAt:  "MAY 8, 2026",
			ReadTime:     "7 MIN",
			Title:        "Porcelain vs. natural stone in Texas homes",
			Excerpt:      "How porosity, UV exposure, Hill Country hard water, and thermal expansion affect limestone, travertine, marble, and kiln-fired porcelain.",
			AuthorName:   "Erin Vasquez",
			AuthorRole:   "Co-Founder & Head of Operations",
			AuthorAvatar: "/images/author-erin.jpg",
			CoverImage:   "/images/blog-materials.jpg",
			Featured:     false,
			Sections: []models.BlogSection{
				{
					Heading: "Porcelain: Zero-sealing resilience for busy spaces",
					Paragraphs: []string{
						"Modern rectified Italian and Spanish porcelain replicates veined marble and warm limestone with <0.5% water absorption. It shrugs off dog claws, patio UV rays, and acidic kitchen spills with minimal maintenance.",
					},
				},
				{
					Heading: "Natural Stone: Warmth and timeless patina",
					Paragraphs: []string{
						"Nothing matches the organic depth of honed travertine or Texas limestone underfoot. Paired with a crack-suppression membrane and breathable penetrating sealer, natural stone wears gracefully for generations.",
					},
				},
			},
		},
		{
			ID:           "post-4",
			Slug:         "regrout-or-replace-how-to-decide",
			Category:     "REPAIR",
			PublishedAt:  "APR 14, 2026",
			ReadTime:     "5 MIN",
			Title:        "Regrout or replace? How to decide",
			Excerpt:      "When diamond-blade regrouting and movement-joint restoration save thousands—and when cracked tiles point to a deeper substrate failure.",
			AuthorName:   "Sam Calloway",
			AuthorRole:   "Founder & Master Setter",
			AuthorAvatar: "/images/author-sam.jpg",
			CoverImage:   "/images/blog-repair.jpg",
			Featured:     false,
			Sections: []models.BlogSection{
				{
					Heading: "When targeted restoration beats tear-out",
					Paragraphs: []string{
						"If your tiles sound solid when tapped and only surface cement grout has stained or crumbled—especially at perimeter expansion joints—removing old grout to full depth and refilling with urethane or epoxy grout restores a floor in a single day.",
					},
				},
				{
					Heading: "When replacement is the only honest fix",
					Paragraphs: []string{
						"If multiple tiles tent upward, rock underfoot, or crack in a continuous jagged line across a room, the underlying mortar bond or slab joint has failed. Regrouting alone will crack again within months.",
					},
				},
			},
		},
	}

	towns := []string{
		"Dripping Springs",
		"Austin",
		"Bee Cave",
		"Lakeway",
		"Wimberley",
		"Fredericksburg",
		"Boerne",
		"Marble Falls",
		"Johnson City",
		"San Marcos",
	}

	return models.ContentBundle{
		Services:         services,
		BeforeAfterCases: beforeAfterCases,
		Testimonials:     testimonials,
		RecentWork:       recentWork,
		HomeFaqs:         homeFaqs,
		FaqCategories:    faqCategories,
		GoogleReviews:    googleReviews,
		GlossaryTerms:    glossaryTerms,
		BlogPosts:        blogPosts,
		Towns:            towns,
	}
}
