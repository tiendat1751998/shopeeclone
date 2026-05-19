package clickhouse

import (
	"context"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn struct {
	pool *pgxpool.Pool
}

func NewConn(ctx context.Context, dsn string) (*Conn, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse ch dsn: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 10 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create ch pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping ch: %w", err)
	}
	return &Conn{pool: pool}, nil
}

func (c *Conn) Close() { c.pool.Close() }

func (c *Conn) InsertViewerEvent(ctx context.Context, roomID, userID, eventType string, timestamp time.Time) error {
	query := `INSERT INTO live_viewer_events (room_id, user_id, event_type, timestamp) VALUES ($1,$2,$3,$4)`
	_, err := c.pool.Exec(ctx, query, roomID, userID, eventType, timestamp)
	return err
}

func (c *Conn) InsertEngagementEvent(ctx context.Context, roomID, userID, eventType string, value int64, timestamp time.Time) error {
	query := `INSERT INTO live_engagement_events (room_id, user_id, event_type, value, timestamp) VALUES ($1,$2,$3,$4,$5)`
	_, err := c.pool.Exec(ctx, query, roomID, userID, eventType, value, timestamp)
	return err
}

func (c *Conn) GetConcurrentViewers(ctx context.Context, roomID string, windowMin int) (int64, error) {
	query := `SELECT COUNT(DISTINCT user_id) FROM live_viewer_events
		WHERE room_id=$1 AND event_type='heartbeat' AND timestamp > now() - INTERVAL $2 MINUTE`
	var count int64
	err := c.pool.QueryRow(ctx, query, roomID, windowMin).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("concurrent viewers: %w", err)
	}
	return count, nil
}

func (c *Conn) GetEngagementSummary(ctx context.Context, roomID string) (likes, gifts, shares int64, err error) {
	query := `SELECT
		COALESCE(SUM(CASE WHEN event_type='like' THEN value END), 0),
		COALESCE(SUM(CASE WHEN event_type='gift' THEN value END), 0),
		COALESCE(SUM(CASE WHEN event_type='share' THEN value END), 0)
		FROM live_engagement_events WHERE room_id=$1`
	err = c.pool.QueryRow(ctx, query, roomID).Scan(&likes, &gifts, &shares)
	return
}
