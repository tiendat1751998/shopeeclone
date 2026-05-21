export interface Product {
  id: string; shop_id: string; name: string; description: string;
  category_id: string; brand: string; status: ProductStatus;
  condition: string; weight: number; dimensions: string;
  metadata?: Record<string, unknown>; version: number;
  created_at: string; updated_at: string;
  skus?: SKU[]; categories?: Category[]; attributes?: ProductAttribute[];
  media?: Media[]; shop?: Shop; rating?: ProductRating; sold_count?: number;
}
export type ProductStatus = "draft"|"pending_moderation"|"active"|"inactive"|"archived"|"rejected"|"suspended";
export interface SKU {
  id: string; product_id: string; sku_code: string; name: string;
  price: number; compare_price: number; currency: string;
  stock: number; reserved_stock: number; weight: number; dimensions: string;
  status: SKUStatus; attributes?: Record<string, string>; sort_order: number;
  created_at: string; updated_at: string;
}
export type SKUStatus = "active"|"inactive"|"sold_out";
export interface Category {
  id: string; parent_id?: string; name: string; slug: string;
  description: string; image_url: string; sort_order: number;
  is_active: boolean; depth: number; path: string;
  children?: Category[]; product_count?: number;
}
export interface ProductAttribute { id: string; name: string; value: string; }
export interface Media {
  id: string; product_id: string; sku_id?: string;
  type: "image"|"video"; url: string; thumbnail_url: string;
  alt_text: string; sort_order: number; status: string;
}
export interface Shop {
  id: string; name: string; avatar_url: string; rating: number;
  follower_count: number; product_count: number;
  response_rate: number; response_time: string; is_official: boolean;
}
export interface ProductRating { average: number; count: number; distribution: Record<number, number>; }
export interface CartItem {
  id: string; product_id: string; sku_id: string; shop_id: string;
  name: string; image_url: string; price: number; original_price?: number;
  currency: string; quantity: number; stock: number; sku_name: string;
  is_selected: boolean; shop_name: string;
}
export interface Cart {
  items: CartItem[]; total_items: number; subtotal: number; currency: string;
}
export interface ShippingAddress {
  id: string; name: string; phone: string;
  address_line1: string; address_line2?: string;
  city: string; state: string; postal_code: string; country: string;
  is_default: boolean;
}
export interface Order {
  id: string; order_number: string; status: OrderStatus;
  subtotal: number; shipping_fee: number; discount: number; total: number;
  currency: string; created_at: string;
}
export type OrderStatus = "pending"|"awaiting_payment"|"paid"|"processing"|"packed"|"shipped"|"delivered"|"completed"|"cancelled"|"refunded";
export interface ApiResponse<T> { success: boolean; data: T; message?: string; code?: string; error?: string; }
export interface LoginRequest { email: string; password: string; }
export interface RegisterRequest { email: string; username: string; password: string; display_name?: string; }
export interface User {
  id: string; email: string; username: string; display_name: string;
  phone: string; avatar_url: string; status: string; created_at: string;
}
export interface AuthTokens { access_token: string; refresh_token: string; expires_in: number; token_type: string; }
export interface PaginatedResponse<T> { data: T[]; total: number; page: number; page_size: number; total_pages: number; }
export interface SearchFilters {
  query?: string; category_id?: string; min_price?: number; max_price?: number;
  brand?: string[]; sort_by?: string; page?: number; page_size?: number;
}
export interface SearchResult {
  products: Product[]; total: number; page: number; page_size: number;
  total_pages: number;
}
