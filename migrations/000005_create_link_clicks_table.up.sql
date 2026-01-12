CREATE TABLE IF NOT EXISTS link_clicks (
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    referer TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_link_clicks_link_date ON link_clicks (link_id, created_at);
CREATE INDEX IF NOT EXISTS idx_link_clicks_link_ip_address ON link_clicks(link_id, ip_address);
