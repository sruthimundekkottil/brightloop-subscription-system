// Package common holds request and response pieces shared by every API.
package common

// PageReq is embedded in every list request.
type PageReq struct {
	Page     int `json:"page"     d:"1"  v:"min:1"`
	PageSize int `json:"pageSize" d:"25" v:"between:1,100"`
}

func (p PageReq) Offset() int { return (p.Page - 1) * p.PageSize }

type PageMeta struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

// List is embedded in every list response (GoFrame wants response names ending in Res):
//
//	type UserListRes struct{ common.List[UserItem] }
//
// middleware.Response renders it as { "data": items, "meta": meta }.
type List[T any] struct {
	Items []T
	Meta  PageMeta
}

func NewList[T any](items []T, p PageReq, total int) List[T] {
	if items == nil {
		items = []T{} // "data": [] rather than null
	}
	return List[T]{items, PageMeta{p.Page, p.PageSize, total}}
}

func (l *List[T]) Paged() (any, PageMeta) { return l.Items, l.Meta }
