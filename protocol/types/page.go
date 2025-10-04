package types

type PageReq struct {
	PageNum  int64 `json:"pageNum,default=1" example:"1"`          // pageNum is the page number
	PageSize int64 `json:"pageSize,default=10" example:"10"`       // pageSize is the page size
	Cursor   int64 `json:"cursor,omitempty,default=0" example:"0"` // cursor is the cursor for the next page
}

func (p PageReq) GetOffset() int64 {
	if p.PageNum == 0 {
		return 0
	}
	return (p.PageNum - 1) * p.PageSize
}

func (p PageReq) GetNextPageOffset() int64 {
	if p.PageNum == 0 {
		p.PageNum = 1
	}
	return p.PageNum * p.PageSize
}

type PageResp struct {
	Total   int64 `json:"total,omitempty" example:"100"`    // the total number of items; 0 means no items or not supported
	HasMore bool  `json:"hasMore,omitempty" example:"true"` // whether there are more items to fetch
}
