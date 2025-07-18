package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"personal-project-core/constant"
	svCtx "personal-project-core/context"
	"personal-project-core/helpers"
	"personal-project-core/models"
)

type UserServer struct {
	sc *svCtx.ServiceContext
}

func NewUserServer(sc *svCtx.ServiceContext) *UserServer {
	return &UserServer{
		sc: sc,
	}
}

func (u *UserServer) PostUsers(ctx context.Context, request PostUsersRequestObject) (PostUsersResponseObject, error) {
	// hash password
	password, err := helpers.HashPassword(request.Body.Password)
	if err != nil {
		u.sc.Log.Error("Failed to hash password", zap.Error(err))
		return PostUsers400JSONResponse{
			Message: "Failed to hash password",
			Code:    constant.PREFIX_USER + constant.ERROR_CREATE_FAILED,
		}, nil
	}

	// check username is existed
	cond := map[string]any{
		"username": request.Body.Username,
	}

	user, _ := u.sc.UserRepo.FindOne(ctx, cond, []string{})
	if user != nil {
		u.sc.Log.Error("Username is duplicated", zap.Any("request", request.Body))
		return PostUsers400JSONResponse{
			Message: "Username is duplicated",
			Code:    constant.PREFIX_USER + constant.ERROR_CREATE_FAILED,
		}, nil
	}
	newUser := models.User{
		Name:     request.Body.Name,
		Username: request.Body.Username,
		Email:    request.Body.Email,
		RoleIds:  request.Body.Roles,
		Password: password,
	}
	fmt.Println("newUser", newUser)
	err = u.sc.UserRepo.Create(ctx, newUser)
	if err != nil {
		u.sc.Log.Error("Failed to create user", zap.Error(err))
		return PostUsers400JSONResponse{
			Message: "Failed to create user",
			Code:    constant.PREFIX_USER + constant.ERROR_CREATE_FAILED,
		}, nil
	}
	return PostUsers201JSONResponse{
		Message: "Created successfully",
	}, nil
}

func (u *UserServer) GetUsers(ctx context.Context, request GetUsersRequestObject) (GetUsersResponseObject, error) {
	limit := request.Params.Limit
	page := request.Params.Page
	offset := (page - 1) * limit
	cond := map[string]any{}
	if request.Params.Keyword != nil {
		cond["OR"] = map[string]any{
			"username LIKE": *request.Params.Keyword,
			"name LIKE":     *request.Params.Keyword,
		}
	}
	if request.Params.IsActive != nil {
		cond["is_active"] = *request.Params.IsActive
	}
	result := u.sc.UserRepo.Paginate(ctx, cond, []string{}, limit, offset, request.Params.Sort)

	users := make([]User, 0)
	for _, user := range result.Items {
		users = append(users, User{
			Id:       user.Id.String(),
			Email:    user.Email,
			IsActive: user.IsActive,
			Name:     user.Name,
			Username: user.Username,
		})
	}
	return GetUsers200JSONResponse{
		Data: &ListResponse{
			Items: users,
			Limit: result.Limit,
			Page:  result.Page,
			Pages: result.Pages,
			Total: int(result.Total),
		},
		Message: nil,
	}, nil
}

func (u *UserServer) GetUsersId(ctx context.Context, request GetUsersIdRequestObject) (GetUsersIdResponseObject, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		u.sc.Log.Error("id is not valid", zap.Error(err), zap.Any("id", request.Id))
		return GetUsersId400JSONResponse{
			Message: "id is not valid",
			Code:    constant.PREFIX_USER + constant.ERROR_INVALID_REQUEST,
		}, nil
	}

	user, err := u.sc.UserRepo.First(ctx, id)
	if err != nil {
		u.sc.Log.Error("failed to get user", zap.Error(err))
		return GetUsersId400JSONResponse{
			Message: "failed to get user",
			Code:    constant.PREFIX_USER + constant.ERROR_DATA_NOT_FOUND,
		}, nil
	}
	var companyId, supplierId string
	if user.CompanyId != nil {
		companyId = user.CompanyId.String()
	}
	if user.SupplierId != nil {
		supplierId = user.SupplierId.String()
	}
	rsp := User{
		Id:         user.Id.String(),
		CompanyId:  &companyId,
		Email:      user.Email,
		IsActive:   user.IsActive,
		Name:       user.Name,
		Stores:     (*[]string)(&user.StoreIds),
		SupplierId: &supplierId,
		Username:   user.Username,
	}
	return GetUsersId200JSONResponse{
		Message: nil,
		Data:    rsp,
	}, nil
}

func (u *UserServer) PutUsersId(ctx context.Context, request PutUsersIdRequestObject) (PutUsersIdResponseObject, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		u.sc.Log.Error("id is not valid", zap.Error(err), zap.Any("id", request.Id))
		return PutUsersId400JSONResponse{
			Message: "id is not valid",
			Code:    constant.PREFIX_USER + constant.ERROR_INVALID_REQUEST,
		}, nil
	}
	password, err := helpers.HashPassword(*request.Body.Password)
	if err != nil {
		u.sc.Log.Error("failed to hash password", zap.Error(err))
		return PutUsersId400JSONResponse{
			Message: "failed to hash password",
			Code:    constant.PREFIX_USER + constant.ERROR_UPDATE_FAILED,
		}, nil
	}
	updateData := map[string]any{
		"email":       *request.Body.Email,
		"name":        *request.Body.Name,
		"is_active":   *request.Body.IsActive,
		"password":    password,
		"roles":       pq.StringArray(*request.Body.Roles),
		"stores":      pq.StringArray(*request.Body.Stores),
		"supplier_id": *request.Body.SupplierId,
	}

	err = u.sc.UserRepo.Update(ctx, id, updateData)
	if err != nil {
		u.sc.Log.Error("failed to update user", zap.Error(err))
		return PutUsersId400JSONResponse{
			Message: "failed to update user",
			Code:    constant.PREFIX_USER + constant.ERROR_UPDATE_FAILED,
		}, nil
	}
	return PutUsersId200JSONResponse{
		Message: "updated successfully",
	}, nil
}

func (u *UserServer) DeleteUsersId(ctx context.Context, request DeleteUsersIdRequestObject) (DeleteUsersIdResponseObject, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		u.sc.Log.Error("id is not valid", zap.Error(err), zap.Any("id", request.Id))
		return DeleteUsersId400JSONResponse{
			Message: "id is not valid",
			Code:    constant.PREFIX_USER + constant.ERROR_INVALID_REQUEST,
		}, nil
	}

	err = u.sc.UserRepo.Delete(ctx, id)
	if err != nil {
		u.sc.Log.Error("failed to delete user", zap.Error(err))
		return DeleteUsersId400JSONResponse{
			Message: "failed to delete user",
			Code:    constant.PREFIX_USER + constant.ERROR_DELETE_FAILED,
		}, nil
	}
	return DeleteUsersId200JSONResponse{
		Message: "deleted successfully",
	}, nil
}
