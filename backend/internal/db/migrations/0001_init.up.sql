CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email varchar(255) NOT NULL UNIQUE,
    password_hash text NOT NULL,
    name varchar(255) NOT NULL,
    onboarding_focus varchar(50),
    preferred_language varchar(5) NOT NULL DEFAULT 'id',
    preferred_currency varchar(5) NOT NULL DEFAULT 'IDR',
    email_verified boolean NOT NULL DEFAULT false,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE refresh_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);

CREATE TABLE email_verification_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE password_reset_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE brands (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL,
    slug varchar(255) NOT NULL UNIQUE,
    logo_url text,
    is_local boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE products (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id uuid NOT NULL REFERENCES brands (id) ON DELETE CASCADE,
    name varchar(500) NOT NULL,
    slug varchar(500) NOT NULL UNIQUE,
    category varchar(100) NOT NULL,
    attributes jsonb NOT NULL DEFAULT '{}'::jsonb,
    is_limited boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_products_brand_id ON products (brand_id);
CREATE INDEX idx_products_category ON products (category);
CREATE INDEX idx_products_is_limited ON products (is_limited) WHERE is_limited = true;
CREATE INDEX idx_products_attributes ON products USING gin (attributes);

CREATE TABLE product_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    locale varchar(5) NOT NULL,
    name varchar(500) NOT NULL,
    description text,
    UNIQUE (product_id, locale)
);

CREATE TABLE product_images (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    url text NOT NULL,
    sort_order int NOT NULL DEFAULT 0
);

CREATE INDEX idx_product_images_product_id ON product_images (product_id);

CREATE TABLE stores (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL,
    type varchar(50) NOT NULL,
    affiliate_network varchar(100),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE product_offers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    store_id uuid NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
    price numeric(12, 2) NOT NULL,
    currency varchar(5) NOT NULL DEFAULT 'IDR',
    in_stock boolean NOT NULL DEFAULT true,
    size varchar(20),
    affiliate_url text NOT NULL,
    scraped_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_product_offers_product_id ON product_offers (product_id);
CREATE INDEX idx_product_offers_store_id ON product_offers (store_id);
CREATE INDEX idx_product_offers_in_stock ON product_offers (in_stock);

CREATE TABLE price_history (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    store_id uuid NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
    price numeric(12, 2) NOT NULL,
    recorded_date date NOT NULL,
    UNIQUE (product_id, store_id, recorded_date)
);

CREATE INDEX idx_price_history_product_id_date ON price_history (product_id, recorded_date);

CREATE TABLE reviews (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    rating int NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment text,
    fit_feedback varchar(20),
    is_flagged boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (product_id, user_id)
);

CREATE INDEX idx_reviews_product_id ON reviews (product_id);

CREATE TABLE review_reports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id uuid NOT NULL REFERENCES reviews (id) ON DELETE CASCADE,
    reported_by uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    reason varchar(255),
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (review_id, reported_by)
);

CREATE TABLE wishlists (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    product_id uuid NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    alert_active boolean NOT NULL DEFAULT false,
    alert_type varchar(50),
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, product_id)
);

CREATE INDEX idx_wishlists_user_id ON wishlists (user_id);

CREATE TABLE notifications (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type varchar(50) NOT NULL,
    title varchar(255) NOT NULL,
    body text NOT NULL,
    action_url text,
    is_read boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_user_id ON notifications (user_id, is_read);

CREATE TABLE user_size_preferences (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    brand_id uuid NOT NULL REFERENCES brands (id) ON DELETE CASCADE,
    size varchar(20) NOT NULL,
    UNIQUE (user_id, brand_id)
);

CREATE TABLE size_conversion_matrix (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    reference_brand_id uuid NOT NULL REFERENCES brands (id) ON DELETE CASCADE,
    target_brand_id uuid NOT NULL REFERENCES brands (id) ON DELETE CASCADE,
    reference_size varchar(20) NOT NULL,
    target_size varchar(20) NOT NULL,
    category varchar(100),
    UNIQUE (reference_brand_id, target_brand_id, reference_size, category)
);

CREATE TABLE exchange_rates (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    base_currency varchar(5) NOT NULL,
    target_currency varchar(5) NOT NULL,
    rate numeric(18, 6) NOT NULL,
    recorded_date date NOT NULL,
    UNIQUE (base_currency, target_currency, recorded_date)
);

CREATE TABLE audit_logs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid REFERENCES users (id) ON DELETE SET NULL,
    action varchar(100) NOT NULL,
    ip_address varchar(64),
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);
