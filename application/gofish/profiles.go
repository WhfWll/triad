package gofish

import (
	"context"
	"smart/api/typespec/gofishtypespec"
	"smart/tools/gophish"
	"time"
)

func (biz *GoPhishBiz) ProfileGetAllBiz(ctx context.Context, req *gofishtypespec.GetListInfoReq) (interface{}, error) {
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
	info, err := client.SendingProfiles().GetAll()
	if err != nil {
		return nil, err
	}

	if len(info) == 0 {
		return map[string]interface{}{
			"sendingProfile": make([]gophish.SendingProfile, 0),
			"total":          0,
		}, nil
	}

	// 使用分页方法
	list, total := paginate(filterItemsBySearch(info, req.Search, biz.sendingProfileSearchFunc), req.Page, req.Size)
	return map[string]interface{}{
		"sendingProfile": list,
		"total":          total,
	}, nil
}

func (biz *GoPhishBiz) ProfileGetDetailBiz(ctx context.Context, id int64) (interface{}, error) {
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
	return client.SendingProfiles().GetByID(id)
}

func (biz *GoPhishBiz) ProfileCreateBiz(ctx context.Context, req *gofishtypespec.CreateProfileReq) (interface{}, error) {
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
	headers := make([]gophish.Header, 0, len(req.Headers))
	for _, h := range req.Headers {
		headers = append(headers, gophish.Header{Key: h.Key, Value: h.Value})
	}
	cReq := &gophish.CreateSendingProfileRequest{
		Name:             req.Name,
		Username:         req.Username,
		Password:         req.Password,
		Host:             req.Host,
		InterfaceType:    req.InterfaceType,
		FromAddress:      req.FromAddress,
		IgnoreCertErrors: req.IgnoreCertErrors,
		Headers:          headers,
	}
	return client.SendingProfiles().Create(cReq)
}

func (biz *GoPhishBiz) ProfileUpdateBiz(ctx context.Context, req *gofishtypespec.UpdateProfileReq) (interface{}, error) {
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
	headers := make([]gophish.Header, 0, len(req.Headers))
	for _, h := range req.Headers {
		headers = append(headers, gophish.Header{Key: h.Key, Value: h.Value})
	}
	up := &gophish.SendingProfile{
		ID:               req.ID,
		Name:             req.Name,
		Username:         req.Username,
		Password:         req.Password,
		Host:             req.Host,
		InterfaceType:    req.InterfaceType,
		FromAddress:      req.FromAddress,
		IgnoreCertErrors: req.IgnoreCertErrors,
		Headers:          headers,
		ModifiedDate:     time.Now().Format(time.RFC3339),
	}
	return client.SendingProfiles().Update(req.ID, up)
}

func (biz *GoPhishBiz) ProfileDeleteBiz(ctx context.Context, id int64) error {
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
	return client.SendingProfiles().Delete(id)
}

func (biz *GoPhishBiz) ProfileSendTestEmailBiz(ctx context.Context, req *gofishtypespec.SendTestEmailReq) error {
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

	headers := make([]gophish.Header, 0, len(req.Smtp.Headers))
	for _, h := range req.Smtp.Headers {
		headers = append(headers, gophish.Header{Key: h.Key, Value: h.Value})
	}

	up := &gophish.SendTestEmailReq{
		Template:  req.Template,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Position:  req.Position,
		Url:       req.Url,
		Page:      req.Page,
		Smtp: gophish.Smtp{
			Name:             req.Smtp.Name,
			FromAddress:      req.Smtp.FromAddress,
			Host:             req.Smtp.Host,
			Username:         req.Smtp.Username,
			Password:         req.Smtp.Password,
			IgnoreCertErrors: req.Smtp.IgnoreCertErrors,
			Headers:          headers,
		},
	}
	client := gophish.New(hc.BaseUri, apiKey)
	return client.SendingProfiles().SendTestEmail(up)
}
