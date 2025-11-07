-- Track basic view metrics per news
CREATE TABLE IF NOT EXISTS news_metrics (
  news_id     BIGINT PRIMARY KEY REFERENCES news(id) ON DELETE CASCADE,
  views       BIGINT NOT NULL DEFAULT 0,
  last_view_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_news_metrics_last_view ON news_metrics(last_view_at DESC);
