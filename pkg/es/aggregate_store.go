package es

import (
	"context"

	"github.com/pkg/errors"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/es/serializer"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/tracing"

	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Load es.Aggregate events using snapshots with given frequency
func (p *pgEventStore) Load(ctx context.Context, aggregate Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "pgEventStore.Load")
	defer span.Finish()
	span.LogFields(log.String("aggregate", aggregate.String()))

	snapshot, err := p.GetSnapshot(ctx, aggregate.GetID())
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return tracing.TraceWithErr(span, err)
	}

	if snapshot != nil {
		if err := serializer.Unmarshal(snapshot.State, aggregate); err != nil {
			p.log.Errorf("(Load) serializer.Unmarshal err: %v", err)
			return tracing.TraceWithErr(span, err)
		}

		err := p.loadAggregateEventsByVersion(ctx, aggregate)
		if err != nil {
			return err
		}

		p.log.Debugf("(Load Aggregate By Version) aggregate: %s", aggregate.String())
		span.LogFields(log.String("aggregate with events", aggregate.String()))
		return nil
	}

	err = p.loadEvents(ctx, aggregate)
	if err != nil {
		return nil
	}

	p.log.Debugf("(Load Aggregate): aggregate: %s", aggregate.String())
	span.LogFields(log.String("aggregae with events", aggregate.String()))
	return nil
}

// Save es.Aggregate events using snapshots with when given frequency
func (p *pgEventStore) Save(ctx context.Context, aggregate Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "pgEventStore.Save")
	defer span.Finish()
	span.LogFields(log.String("aggregate", aggregate.String()))

	// Check if any event not create then dont need save any thing
	if len(aggregate.GetChanges()) == 0 {
		p.log.Debug("(Save) aggregate.GetChanges()) == 0")
		span.LogFields(log.Int("events", len(aggregate.GetChanges())))
		return nil
	}

	// Save process include save event and create snapshot with implement logs and process event then used transaction
	// Begin and create transaction
	tx, err := p.db.Begin(ctx)
	if err != nil {
		p.log.Errorf("(Save) db.Begin err: %v", err)
		return tracing.TraceWithErr(span, errors.Wrap(err, "db.Begin"))
	}

	// this functionality run after all process
	defer func() {
		if tx != nil {
			if txErr := tx.Rollback(ctx); txErr != nil && !errors.Is(txErr, pgx.ErrTxClosed) {
				err = txErr
				tracing.TraceWithErr(span, err)
				return
			}
		}
	}()

	// consider created events
	changes := aggregate.GetChanges()
	events := make([]Event, 0, len(changes))

	// serialize all event with error tracing
	for i := range changes {
		event, err := p.serializer.SerializeEvent(aggregate, changes[i])
		if err != nil {
			p.log.Errorf("(Save) serializer.SerializeEvent err: %v", err)
			return tracing.TraceWithErr(span, errors.Wrap(err, "serializer.SerializeEvent"))
		}

		events = append(events, event)
	}

	// save event with transaction and error tracing
	if err := p.saveEventsTx(ctx, tx, events); err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, "saveEventTx"))
	}

	// consider aggregate version with snapshotFrequency for process Snapshot and save snapshot with transaction
	if aggregate.GetVersion()%p.cfg.SnapshotFrequency == 0 {
		aggregate.ToSnapshot()
		if err := p.saveSnapshotTx(ctx, tx, aggregate); err != nil {
			return tracing.TraceWithErr(span, errors.Wrap(err, "saveSnapshotTx"))
		}
	}

	// run processEvents
	if err := p.processEvents(ctx, events); err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, "processEvents"))
	}

	// trace process and commit transaction
	p.log.Debugf("(Save Aggregate): aggregate: %s", aggregate.String())
	span.LogFields(log.String("aggregate with events", aggregate.String()))
	return tx.Commit(ctx)
}
