package platform

import "context"

// Actor is the logged-in user behind a request. The zero value means a system job.
type Actor struct {
	UserID     uint64
	Role       string // admin, staff, customer
	CustomerID uint64 // 0 unless Role is customer
}

type actorKey struct{}

func WithActor(ctx context.Context, a Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, a)
}

func ActorFrom(ctx context.Context) Actor {
	a, _ := ctx.Value(actorKey{}).(Actor)
	return a
}
