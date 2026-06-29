package users

type FindByIdParams struct {
	Id int `uri:"id" binding:"required"`
}

type OrderSort string

const (
	OrderAsc  OrderSort = "asc"
	OrderDesc OrderSort = "desc"
)

type SortBy string

const (
	SortByCreatedAt SortBy = "createdAt"
	SortByName      SortBy = "name"
	SortByEmail     SortBy = "email"
)

type PaginationQuery struct {
	Page      int       `form:"page,default=1" binding:"min=1"`
	Limit     int       `form:"limit,default=10" binding:"min=1,max=500"`
	SortBy    SortBy    `form:"sortBy,default=createdAt" binding:"oneof=createdAt name email"`
	OrderSort OrderSort `form:"orderSort,default=asc" binding:"oneof=asc desc"`
}

type AdminUsersPageRequest struct {
	UserFilter
	PaginationQuery
}
