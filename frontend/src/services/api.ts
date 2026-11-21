import type {
  BlogPost,
  BookingAvailabilityResponse,
  BookingRequestPayload,
  BookingResponse,
  ContentBundle,
  EstimateRequestPayload,
  EstimateResponse,
  ServiceItem,
} from '@/types/copely';
import { assetUrl } from '@/utils/assetUrl';

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? '').trim().replace(/\/+$/, '');

interface ApiErrorPayload {
  success?: boolean;
  error?: string;
}

export class ApiError extends Error {
  public readonly status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

function apiUrl(path: string): string {
  return API_BASE_URL ? `${API_BASE_URL}${path}` : path;
}

function normalizeServiceItem(item: ServiceItem): ServiceItem {
  return { ...item, imageUrl: assetUrl(item.imageUrl) };
}

function normalizeBlogPost(post: BlogPost): BlogPost {
  return {
    ...post,
    authorAvatar: assetUrl(post.authorAvatar),
    coverImage: assetUrl(post.coverImage),
  };
}

function normalizeContentBundle(bundle: ContentBundle): ContentBundle {
  return {
    ...bundle,
    services: bundle.services.map(normalizeServiceItem),
    beforeAfterCases: bundle.beforeAfterCases.map((item) => ({
      ...item,
      beforeImage: assetUrl(item.beforeImage),
      afterImage: assetUrl(item.afterImage),
    })),
    testimonials: bundle.testimonials.map((item) => ({
      ...item,
      avatarUrl: assetUrl(item.avatarUrl),
      heroImageUrl: assetUrl(item.heroImageUrl),
    })),
    recentWork: bundle.recentWork.map((item) => ({
      ...item,
      imageUrl: assetUrl(item.imageUrl),
    })),
    blogPosts: bundle.blogPosts.map(normalizeBlogPost),
  };
}

async function requestJson<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(apiUrl(path), {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      ...(init?.headers ?? {}),
    },
  });

  if (!response.ok) {
    let message = `Request failed (${response.status})`;
    let parsedJsonError = false;
    try {
      const errBody = (await response.json()) as ApiErrorPayload;
      if (errBody.error) {
        message = errBody.error;
        parsedJsonError = true;
      }
    } catch {
      // Fall back to HTTP status text when error body is non-JSON.
    }
    if (
      !parsedJsonError &&
      !API_BASE_URL &&
      init?.method === 'POST' &&
      (response.status === 404 || response.status === 405)
    ) {
      message =
        'Online booking and estimate submissions require a live API backend. Please call (512) 555-0142.';
    }
    throw new ApiError(message, response.status);
  }

  return (await response.json()) as T;
}

export async function fetchContent(): Promise<ContentBundle> {
  return normalizeContentBundle(await requestJson<ContentBundle>('/api/content'));
}

export async function fetchServiceBySlug(slug: string): Promise<ServiceItem> {
  return normalizeServiceItem(
    await requestJson<ServiceItem>(`/api/services/${encodeURIComponent(slug)}`),
  );
}

export async function fetchPostBySlug(slug: string): Promise<BlogPost> {
  return normalizeBlogPost(
    await requestJson<BlogPost>(`/api/posts/${encodeURIComponent(slug)}`),
  );
}

export function submitEstimate(payload: EstimateRequestPayload): Promise<EstimateResponse> {
  return requestJson<EstimateResponse>('/api/estimates', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function fetchBookingAvailability(month: string): Promise<BookingAvailabilityResponse> {
  return requestJson<BookingAvailabilityResponse>(
    `/api/bookings/availability?month=${encodeURIComponent(month)}`,
  );
}

export function createBooking(payload: BookingRequestPayload): Promise<BookingResponse> {
  return requestJson<BookingResponse>('/api/bookings', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
