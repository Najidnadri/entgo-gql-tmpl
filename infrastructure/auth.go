package infrastructure

// import (
// 	"context"
// 	"shopolah/ent"
// 	"shopolah/ent/user"
// )

// type contextKey struct {
// 	name string
// }

// var UserCtxKey = &contextKey{"viewer"}

// type ViewerObj struct {
// 	User          *ent.User
// 	FirebaseXid   *string
// 	EmailVerified *bool
// 	Role          user.Role
// 	UserType      UserType
// }

// type Viewer interface {
// 	IsValid() bool
// 	IsAdmin() bool
// 	IsEmployee() bool
// }

// var _ Viewer = &ViewerObj{}

// func (v *ViewerObj) IsValid() bool {
// 	if v != nil && v.User != nil {
// 		return true
// 	}
// 	return false
// }

// func (v *ViewerObj) IsAdmin() bool {
// 	if v.IsValid() && v.Role == user.RoleAdmin && v.UserType == UserTypeAdmin {
// 		return true
// 	}
// 	return false
// }

// func (v *ViewerObj) IsEmployee() bool {
// 	if v.IsValid() && (v.Role == user.RoleEmployee || v.Role == user.RoleAdmin) {
// 		return true
// 	}
// 	return false
// }

// func (v *ViewerObj) IsCustomer() bool {
// 	if v.IsValid() && v.UserType == UserTypeCustomer {
// 		return true
// 	}
// 	return false
// }

// func (v *ViewerObj) IsSeller() bool {
// 	if v.IsValid() && v.UserType == UserTypeSeller {
// 		return true
// 	}
// 	return false
// }

// func AuthFromContext(ctx context.Context) *ViewerObj {
// 	raw, _ := ctx.Value(UserCtxKey).(*ViewerObj)
// 	return raw
// }

// type UserType string

// const (
// 	UserTypeCustomer UserType = "customer"
// 	UserTypeSeller   UserType = "seller"
// 	UserTypeAdmin    UserType = "admin"
// )
