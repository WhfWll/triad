// Package application
// @Author bcy2007  2025/12/23 11:05
package application

import (
	"context"
	"smart/api/typespec"
	"smart/tools/facade"
)

type ReverseServer struct{}

func (r *ReverseServer) Start(_ context.Context, req *typespec.ReverseServerStartReq) {
	facade.NewFacadeServer(req.Host, req.Port, 1800)
}

func (r *ReverseServer) Status(_ context.Context, resp *typespec.ReverseServerStatusResp) {
	status, info := facade.GetFacadeServerStatus()
	resp.Status = status
	if status {
		resp.ReverseUrl = info
	} else {
		resp.ErrInfo = info
	}
}

func (r *ReverseServer) Stop(_ context.Context) {
	facade.CloseFacadeServer()
}

func (r *ReverseServer) ClearMessage(_ context.Context) {
	facade.ClearMessage()
}

func (r *ReverseServer) Messages(_ context.Context, req *typespec.ReverseServerMessageReq, resp *typespec.ReverseServerMessageResp) {
	messages, total := facade.ReadMessage(req.Page, req.Size)
	resp.Total = total
	for _, message := range messages {
		resp.List = append(resp.List, typespec.ReverseServerMessageItem{
			ReverseType: message.FacadeType,
			RemoteAddr:  message.RemoteAddr,
			Token:       message.Token,
			Response:    message.Response,
		})
	}
}
