package gofish

import (
	"context"
	"smart/api/typespec/gofishtypespec"
	"smart/tools/gophish"
	"time"
)

// GetGroupAllBiz .
func (biz *GoPhishBiz) GetGroupAllBiz(ctx context.Context, req *gofishtypespec.GetListInfoReq) (interface{}, error) {
	userId := ctx.Value("uid").(int)
	username := ctx.Value("username").(string)

	apiKey, err := biz.GetGoPhishUserAPiKeyInfo(ctx, userId, username)
	if err != nil {
		return nil, err
	}

	hc, err := getGOPhishServiceInfo()
	if err != nil {
		return nil, err
	}

	client := gophish.New(hc.BaseUri, apiKey)
	info, err := client.Groups().GetAll()
	if err != nil {
		return nil, err
	}

	if info.Total == 0 {
		return info, nil
	}

	groups := filterItemsBySearch(info.Groups, req.Search, biz.groupSearchFunc)
	// 使用分页方法
	list, total := paginate(groups, req.Page, req.Size)
	return gophish.GroupResponse{
		Total:  total,
		Groups: list,
	}, nil
}

// GetGroupDetailBiz .
func (biz *GoPhishBiz) GetGroupDetailBiz(ctx context.Context, req *gofishtypespec.GroupDetailReq) (interface{}, error) {
	userId := ctx.Value("uid").(int)
	username := ctx.Value("username").(string)

	apiKey, err := biz.GetGoPhishUserAPiKeyInfo(ctx, userId, username)
	if err != nil {
		return nil, err
	}

	hc, err := getGOPhishServiceInfo()
	if err != nil {
		return nil, err
	}

	client := gophish.New(hc.BaseUri, apiKey)
	info, err := client.Groups().GetByID(req.ID)

	return info, err
}

// GroupCreateBiz .
func (biz *GoPhishBiz) GroupCreateBiz(ctx context.Context, req *gofishtypespec.CreateGroupReq) (interface{}, error) {
	userId := ctx.Value("uid").(int)
	username := ctx.Value("username").(string)

	apiKey, err := biz.GetGoPhishUserAPiKeyInfo(ctx, userId, username)
	if err != nil {
		return nil, err
	}

	targets := make([]gophish.Target, 0, len(req.Targets))
	for _, target := range req.Targets {
		targets = append(targets, gophish.Target{
			Email:     target.Email,
			FirstName: target.FirstName,
			Position:  target.Position,
			LastName:  target.LastName,
		})
	}

	reqs := gophish.CreateGroupRequest{
		Name:    req.Name,
		Targets: targets,
	}

	hc, err := getGOPhishServiceInfo()
	if err != nil {
		return nil, err
	}

	client := gophish.New(hc.BaseUri, apiKey)
	info, err := client.Groups().Create(&reqs)

	return info, err
}

// GroupUpdateBiz .
func (biz *GoPhishBiz) GroupUpdateBiz(ctx context.Context, req *gofishtypespec.UpdateGroupReq) (interface{}, error) {
	userId := ctx.Value("uid").(int)
	username := ctx.Value("username").(string)

	apiKey, err := biz.GetGoPhishUserAPiKeyInfo(ctx, userId, username)
	if err != nil {
		return nil, err
	}

	targets := make([]gophish.Target, 0, len(req.Targets))
	for _, target := range req.Targets {
		targets = append(targets, gophish.Target{
			Email:     target.Email,
			FirstName: target.FirstName,
			Position:  target.Position,
			LastName:  target.LastName,
		})
	}

	up := &gophish.Group{
		ID:           req.ID,
		Name:         req.Name,
		Targets:      targets,
		ModifiedDate: time.Now().Format(time.RFC3339),
	}

	hc, err := getGOPhishServiceInfo()
	if err != nil {
		return nil, err
	}

	client := gophish.New(hc.BaseUri, apiKey)
	info, err := client.Groups().Update(req.ID, up)

	return info, err
}

// GroupDeleteBiz .
func (biz *GoPhishBiz) GroupDeleteBiz(ctx context.Context, req *gofishtypespec.GroupDetailReq) error {
	userId := ctx.Value("uid").(int)
	username := ctx.Value("username").(string)

	apiKey, err := biz.GetGoPhishUserAPiKeyInfo(ctx, userId, username)
	if err != nil {
		return err
	}

	hc, err := getGOPhishServiceInfo()
	if err != nil {
		return err
	}

	client := gophish.New(hc.BaseUri, apiKey)
	err = client.Groups().Delete(req.ID)

	return err
}
