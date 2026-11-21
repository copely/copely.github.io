export type TileServiceCategory =
  | 'porcelain-ceramic'
  | 'natural-stone'
  | 'wet-rooms'
  | 'kitchen-backsplashes'
  | 'outdoor-patios';

export type ProjectTimeline = 'this-week' | 'next-2-weeks' | 'planning-ahead';

export interface EstimateRequestPayload {
  firstName?: string;
  lastName?: string;
  phone?: string;
  email?: string;
  address?: string;
  town?: string;
  service: TileServiceCategory;
  timeline: ProjectTimeline;
  notes?: string;
}

export interface EstimateResponse {
  success: boolean;
  estimateId: string;
  responseWindow: string;
  service: TileServiceCategory;
  timeline: ProjectTimeline;
  createdAt: string;
}

export interface BookingRequestPayload {
  date: string; // ISO YYYY-MM-DD
  timeSlot: string; // e.g., "08:00"
  town: string;
  service?: TileServiceCategory;
  timeline?: ProjectTimeline;
  customerName?: string;
  customerPhone?: string;
  customerEmail?: string;
  notes?: string;
}

export interface BookingResponse {
  success: boolean;
  bookingId: string;
  date: string;
  timeSlot: string;
  town: string;
  service?: TileServiceCategory;
  timeline?: ProjectTimeline;
  createdAt: string;
}

export interface TimeSlotAvailability {
  time: string;
  available: boolean;
}

export interface DayAvailability {
  date: string;
  dayOfWeek: number;
  available: boolean;
  slots: TimeSlotAvailability[];
}

export interface BookingAvailabilityResponse {
  month: string;
  days: DayAvailability[];
}

export interface ServiceItem {
  id: string;
  slug: TileServiceCategory;
  title: string;
  shortSubtitle: string;
  description: string;
  imageUrl: string;
  warrantyYears: number;
  timelineLabel: string;
  idealFor: string;
  specs: string[];
  processHighlights: string[];
}

export interface BeforeAfterCase {
  id: string;
  title: string;
  subtitle: string;
  location: string;
  beforeImage: string;
  afterImage: string;
  featured: boolean;
}

export interface TestimonialMetric {
  label: string;
  value: string;
}

export interface Testimonial {
  id: string;
  author: string;
  projectType: string;
  location: string;
  avatarUrl: string;
  quote: string;
  heroImageUrl: string;
  metrics: TestimonialMetric[];
}

export interface RecentWorkItem {
  id: string;
  title: string;
  location: string;
  category: string;
  imageUrl: string;
}

export interface FaqItem {
  id: string;
  category: string;
  question: string;
  answer: string;
}

export interface FaqCategoryGroup {
  id: string;
  title: string;
  items: FaqItem[];
}

export interface GoogleReview {
  id: string;
  author: string;
  initials: string;
  rating: number;
  dateLabel: string;
  projectTag: string;
  body: string;
}

export interface GlossaryTerm {
  id: string;
  term: string;
  definition: string;
}

export interface BlogSection {
  heading: string;
  paragraphs: string[];
}

export interface BlogPost {
  id: string;
  slug: string;
  category: string;
  publishedAt: string;
  readTime: string;
  title: string;
  excerpt: string;
  authorName: string;
  authorRole: string;
  authorAvatar: string;
  coverImage: string;
  featured: boolean;
  sections: BlogSection[];
}

export interface ContentBundle {
  services: ServiceItem[];
  beforeAfterCases: BeforeAfterCase[];
  testimonials: Testimonial[];
  recentWork: RecentWorkItem[];
  homeFaqs: FaqItem[];
  faqCategories: FaqCategoryGroup[];
  googleReviews: GoogleReview[];
  glossaryTerms: GlossaryTerm[];
  blogPosts: BlogPost[];
  towns: string[];
}
