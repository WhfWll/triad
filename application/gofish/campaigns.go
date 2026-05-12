package gofish

import (
	"context"
	"fmt"
	"smart/api/typespec/gofishtypespec"
	"smart/tools/gophish"
	"smart/tools/time"
	"strings"
)

// Campaigns
func (biz *GoPhishBiz) CampaignGetAllBiz(ctx context.Context, req *gofishtypespec.CampaignListReq) (interface{}, error) {
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
	info, err := client.Campaigns().GetAll()
	if err != nil {
		return nil, err
	}

	if info.Total == 0 {
		return info, nil
	}

	groups := filterItemsBySearch(info.Campaigns, req.Search, biz.campaignSearchFunc)
	groups = filterItemsByStatus(groups, req.ActiveStatus, biz.campaignStatusFunc)

	// 使用分页方法
	list, total := paginate(groups, req.Page, req.Size)

	return gophish.CampaignsResp{
		Total:     int(total),
		Campaigns: list,
	}, nil
}

func (biz *GoPhishBiz) CampaignGetDetailBiz(ctx context.Context, id int64) (interface{}, error) {
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
	return client.Campaigns().GetByID(id)
}

func (biz *GoPhishBiz) CampaignCreateBiz(ctx context.Context, req *gofishtypespec.CreateCampaignReq) (interface{}, error) {
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
	launchDate := req.LaunchDate
	if launchDate == "" {
		launchDate = time.GenerateCurrentRFC3339WithOffset()
	}
	cReq := &gophish.CreateCampaignRequest{
		Name:       req.Name,
		Template:   gophish.NameInfo{Name: req.Template.Name},
		Url:        req.Url,
		Page:       gophish.NameInfo{Name: req.Page.Name},
		Smtp:       gophish.NameInfo{Name: req.Smtp.Name},
		LaunchDate: launchDate,
		Groups:     req.Groups,
	}

	if req.SendByDate != nil && req.SendByDate != "" {
		cReq.SendByDate = req.SendByDate
	}
	resp, err := client.Campaigns().Create(cReq)
	if err != nil {
		if strings.Contains(err.Error(), "The launch date must be before the \"send emails by\" date") {
			return nil, fmt.Errorf("启动时间必须早于发送截止时间")
		}
		return nil, err
	}
	return resp, nil
}

func (biz *GoPhishBiz) CampaignLaunchBiz(ctx context.Context, id int64) error {
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
	return client.Campaigns().Launch(id)
}

func (biz *GoPhishBiz) CampaignCompleteBiz(ctx context.Context, id int64) error {
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
	return client.Campaigns().Complete(id)
}

func (biz *GoPhishBiz) CampaignUpdateBiz(ctx context.Context, req *gofishtypespec.UpdateCampaignReq) (interface{}, error) {
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
	up := &gophish.Campaign{ID: req.ID, Name: req.Name, Status: req.Status}
	return client.Campaigns().Update(req.ID, up)
}

func (biz *GoPhishBiz) CampaignDeleteBiz(ctx context.Context, id int64) error {
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
	return client.Campaigns().Delete(id)
}

func (biz *GoPhishBiz) CampaignResultBiz(ctx context.Context, id int64) (interface{}, error) {
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
	return client.Campaigns().Results(id)
}
