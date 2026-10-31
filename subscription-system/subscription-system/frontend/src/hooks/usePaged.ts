import { useCallback, useEffect, useState } from 'react'
import type { TablePaginationConfig } from 'antd'
import type { Page, PageMeta } from '../requests/client'

// usePaged loads one page of a list endpoint and gives AntD Table its pagination props.
// Changing the filters goes back to page 1.
export function usePaged<T, F extends object>(fetch: (p: F & { page: number; pageSize: number }) => Promise<Page<T>>, filters: F) {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(25)
  const [rows, setRows] = useState<T[]>([])
  const [meta, setMeta] = useState<PageMeta>({ page: 1, pageSize: 25, total: 0 })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(false)
  const [tick, setTick] = useState(0)
  const filterKey = JSON.stringify(filters)

  useEffect(() => setPage(1), [filterKey])

  useEffect(() => {
    let live = true
    setLoading(true)
    fetch({ ...filters, page, pageSize })
      .then((r) => {
        if (!live) return
        setRows(r.data)
        setMeta(r.meta)
        setError(false)
      })
      .catch(() => live && setError(true))
      .finally(() => live && setLoading(false))
    return () => {
      live = false
    }
    // fetch is a stable module function; filters are compared by value through filterKey
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterKey, page, pageSize, tick])

  const reload = useCallback(() => setTick((t) => t + 1), [])
  const pagination: TablePaginationConfig = {
    current: meta.page,
    pageSize: meta.pageSize,
    total: meta.total,
    showSizeChanger: true,
    pageSizeOptions: [10, 25, 50, 100],
    onChange: (p, s) => {
      setPage(s !== pageSize ? 1 : p)
      setPageSize(s)
    },
  }
  return { rows, loading, error, reload, pagination }
}
