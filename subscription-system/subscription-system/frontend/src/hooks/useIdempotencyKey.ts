import { useMemo } from 'react'

// useIdempotencyKey returns a fresh key each time a form opens and keeps it while the form stays open,
// so a double-clicked or retried submit sends the same key and the server records it once.
// Usage: const key = useIdempotencyKey(open); send('post', url, body, { 'Idempotency-Key': key })
export function useIdempotencyKey(open: boolean): string {
  // eslint-disable-next-line react-hooks/exhaustive-deps
  return useMemo(() => crypto.randomUUID(), [open])
}
