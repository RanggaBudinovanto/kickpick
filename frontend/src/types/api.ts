export interface Product {
  id: string;
  slug: string;
  name: string;
  category: string;
  brand_id: string;
  is_limited: boolean;
  attributes: Record<string, unknown>;
  brand_name: string;
  brand_slug: string;
  min_price: number;
  max_price: number;
  currency: string;
  rating: number;
  image_url: string;
  drop_percent?: number;
}

export interface Offer {
  id: string;
  store_name: string;
  store_type: string;
  price: number;
  currency: string;
  in_stock: boolean;
  size?: string;
  affiliate_url: string;
}

export interface Review {
  id: string;
  rating: number;
  comment: string;
  fit_feedback: string;
  user_name: string;
  user_id?: string;
  created_at: string;
}

export interface ProductDetail {
  product: Product;
  offers: Offer[];
  reviews: Review[];
}

export interface Brand {
  id: string;
  name: string;
  slug: string;
  logo_url: string;
  is_local: boolean;
}

export interface PricePoint {
  date: string;
  price: number;
}

export interface WishlistItem {
  id: string;
  product_id: string;
  product_name: string;
  product_slug: string;
  alert_active: boolean;
  alert_type: string;
}

export interface NotificationItem {
  id: string;
  type: string;
  title: string;
  body: string;
  action_url: string;
  is_read: boolean;
  created_at: string;
}

export interface UserProfile {
  id: string;
  email: string;
  name: string;
  onboarding_focus: string;
  preferred_language: string;
  preferred_currency: string;
  email_verified: boolean;
}
