CREATE TABLE product_views (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    viewed_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_product_views_product_id_viewed_at ON product_views (product_id, viewed_at);
