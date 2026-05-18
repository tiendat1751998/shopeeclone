package mysql

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/services/auth/internal/config"
	"github.com/shopee-clone/shopee/services/auth/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type AuditRepository struct {
	db          *sqlx.DB
	cfg         config.AuditConfig
	mu          sync.Mutex
	buffer      []*domain.AuditLog
	lastFlush   time.Time
	stopCh      chan struct{}
}

func NewAuditRepository(db *sqlx.DB, cfg config.AuditConfig) *AuditRepository {
	r := &AuditRepository{
		db:        db,
		cfg:       cfg,
		buffer:    make([]*domain.AuditLog, 0, cfg.BatchSize),
		lastFlush: time.Now(),
		stopCh:    make(chan struct{}),
	}
	go r.flushLoop()
	return r
}

func (r *AuditRepository) Log(ctx context.Context, log *domain.AuditLog) {
	if !r.cfg.Enabled {
		return
	}

	r.mu.Lock()
	r.buffer = append(r.buffer, log)
	shouldFlush := len(r.buffer) >= r.cfg.BatchSize || time.Since(r.lastFlush) > r.cfg.FlushInterval
	r.mu.Unlock()

	if shouldFlush {
		go r.Flush()
	}
}

func (r *AuditRepository) Flush() {
	r.mu.Lock()
	if len(r.buffer) == 0 {
		r.mu.Unlock()
		return
	}
	batch := r.buffer
	r.buffer = make([]*domain.AuditLog, 0, r.cfg.BatchSize)
	r.lastFlush = time.Now()
	r.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.batchInsert(ctx, batch); err != nil {
		observability.GetLogger().Error("audit flush failed",
			zap.Int("count", len(batch)),
			zap.Error(err),
		)
		r.mu.Lock()
		r.buffer = append(batch, r.buffer...)
		r.mu.Unlock()
	}
}

func (r *AuditRepository) batchInsert(ctx context.Context, logs []*domain.AuditLog) error {
	query := `INSERT INTO audit_logs (id, trace_id, user_id, action, ip, device_id, user_agent, resource, status, detail, created_at)
		VALUES (:id, :trace_id, :user_id, :action, :ip, :device_id, :user_agent, :resource, :status, :detail, :created_at)`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("audit tx begin: %w", err)
	}
	defer tx.Rollback()

	for _, log := range logs {
		_, err := tx.NamedExecContext(ctx, query, log)
		if err != nil {
			return fmt.Errorf("audit insert: %w", err)
		}
	}

	return tx.Commit()
}

func (r *AuditRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.AuditLog, error) {
	ctx, span := otel.Tracer("shopee-auth").Start(ctx, "mysql.audit.find_by_user")
	defer span.End()

	var logs []*domain.AuditLog
	query := `SELECT id, trace_id, user_id, action, ip, device_id, user_agent, resource, status, detail, created_at
		FROM audit_logs WHERE user_id = ?
		ORDER BY created_at DESC LIMIT ? OFFSET ?`

	err := r.db.SelectContext(ctx, &logs, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	span.SetAttributes(attribute.Int("count", len(logs)))
	return logs, nil
}

func (r *AuditRepository) FindByAction(ctx context.Context, action domain.AuditAction, limit, offset int) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	query := `SELECT id, trace_id, user_id, action, ip, device_id, user_agent, resource, status, detail, created_at
		FROM audit_logs WHERE action = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

	err := r.db.SelectContext(ctx, &logs, query, string(action), limit, offset)
	return logs, err
}

func (r *AuditRepository) DeleteOlderThan(ctx context.Context, ttl time.Duration) error {
	cutoff := time.Now().Add(-ttl)
	_, err := r.db.ExecContext(ctx, `DELETE FROM audit_logs WHERE created_at < ?`, cutoff)
	return err
}

func (r *AuditRepository) flushLoop() {
	ticker := time.NewTicker(r.cfg.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.Flush()
		case <-r.stopCh:
			r.Flush()
			return
		}
	}
}

func (r *AuditRepository) Stop() {
	close(r.stopCh)
}
