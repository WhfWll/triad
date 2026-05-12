package gofish

import (
	"context"
	"smart/api/typespec/gofishtypespec"
	"smart/tools/gophish"
)

func (biz *GoPhishBiz) PageGetAllBiz(ctx context.Context, req *gofishtypespec.GetListInfoReq) (interface{}, error) {
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
	info, err := client.LandingPages().GetAll()
	if err != nil {
		return nil, err
	}

	if len(info) == 0 {
		return map[string]interface{}{
			"landingPage": make([]gophish.LandingPage, 0),
			"total":       0,
		}, nil
	}

	// 使用分页方法
	list, total := paginate(filterItemsBySearch(info, req.Search, biz.landingPageSearchFunc), req.Page, req.Size)
	return map[string]interface{}{
		"landingPage": list,
		"total":       total,
	}, nil
}

func (biz *GoPhishBiz) PageGetDetailBiz(ctx context.Context, id int64) (interface{}, error) {
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
	return client.LandingPages().GetByID(id)
}

func (biz *GoPhishBiz) PageCreateBiz(ctx context.Context, req *gofishtypespec.CreatePageReq) (interface{}, error) {
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
	cReq := &gophish.CreateLandingPageRequest{
		Name:               req.Name,
		HTML:               req.HTML,
		CaptureCredentials: req.CaptureCredentials,
		CapturePasswords:   req.CapturePasswords,
		RedirectURL:        req.RedirectURL,
	}
	return client.LandingPages().Create(cReq)
}

func (biz *GoPhishBiz) PageUpdateBiz(ctx context.Context, req *gofishtypespec.UpdatePageReq) (interface{}, error) {
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
	up := &gophish.LandingPage{
		ID:                 req.ID,
		Name:               req.Name,
		HTML:               req.HTML,
		CaptureCredentials: req.CaptureCredentials,
		CapturePasswords:   req.CapturePasswords,
		RedirectURL:        req.RedirectURL,
	}
	return client.LandingPages().Update(req.ID, up)
}

func (biz *GoPhishBiz) PageDeleteBiz(ctx context.Context, id int64) error {
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
	return client.LandingPages().Delete(id)
}
func (biz *GoPhishBiz) PageImportSiteBiz(ctx context.Context, req *gofishtypespec.ImportSiteReq) (interface{}, error) {
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
	cReq := &gophish.ImportSiteRequest{URL: req.URL, IncludeResources: req.IncludeResources}
	return client.LandingPages().ImportSite(cReq)
}
