package gofish

import (
	"context"
	"smart/api/typespec/gofishtypespec"
	"smart/tools/gophish"
)

func (biz *GoPhishBiz) TemplateGetAllBiz(ctx context.Context, req *gofishtypespec.GetListInfoReq) (interface{}, error) {
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
	info, err := client.Templates().GetAll()
	if err != nil {
		return nil, err
	}

	if len(info) == 0 {
		return map[string]interface{}{
			"templates": make([]gophish.Template, 0),
			"total":     0,
		}, nil
	}

	// 使用分页方法
	list, total := paginate(filterItemsBySearch(info, req.Search, biz.templatesSearchFunc), req.Page, req.Size)
	return map[string]interface{}{
		"templates": list,
		"total":     total,
	}, nil
}

func (biz *GoPhishBiz) TemplateGetDetailBiz(ctx context.Context, id int64) (interface{}, error) {
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
	return client.Templates().GetByID(id)
}

func (biz *GoPhishBiz) TemplateCreateBiz(ctx context.Context, req *gofishtypespec.CreateTemplateReq) (interface{}, error) {
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

	attachments := make([]gophish.Attachment, 0)
	for _, val := range req.Attachments {
		attachments = append(attachments, gophish.Attachment{
			Name:    val.Name,
			Type:    val.Type,
			Content: val.Content,
		})
	}
	cReq := &gophish.CreateTemplateRequest{
		Name: req.Name, Subject: req.Subject, Text: req.Text, HTML: req.HTML, EnvelopeSender: req.EnvelopeSender, Attachments: attachments,
	}

	return client.Templates().Create(cReq)
}

func (biz *GoPhishBiz) TemplateUpdateBiz(ctx context.Context, req *gofishtypespec.UpdateTemplateReq) (interface{}, error) {
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
	attachments := make([]gophish.Attachment, 0)
	for _, val := range req.Attachments {
		attachments = append(attachments, gophish.Attachment{
			Name:    val.Name,
			Type:    val.Type,
			Content: val.Content,
		})
	}
	up := &gophish.Template{ID: req.ID, Name: req.Name, Subject: req.Subject, Text: req.Text, HTML: req.HTML, EnvelopeSender: req.EnvelopeSender, Attachments: attachments}
	return client.Templates().Update(req.ID, up)
}

func (biz *GoPhishBiz) TemplateDeleteBiz(ctx context.Context, id int64) error {
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
	return client.Templates().Delete(id)
}

func (biz *GoPhishBiz) TemplateImportEmailBiz(ctx context.Context, req *gofishtypespec.ImportEmailReq) (interface{}, error) {
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
	cReq := &gophish.ImportEmailRequest{Content: req.Content, ConvertLinks: req.ConvertLinks}
	return client.Templates().ImportEmail(cReq)
}
