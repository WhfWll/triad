// Package crons
// @Author bcy2007  2024/12/18 15:18
package crons

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"smart/api/typespec"
	"smart/application"
	"smart/models/mysqls"
	"smart/tools/enums"
)

func RemoteSessionStatusCheck() {
	var (
		ctx           = context.Background()
		remoteSession mysqls.RemoteSession
		vulApp        application.VulEvidenceListApp
	)
	sessions, err := remoteSession.All(ctx, fmt.Sprintf("status = %v", enums.SessionStatusSucc))
	if err != nil {
		log.Errorf("RemoteSessionStatusCheck Error: get all remote session list error: %v", err)
		return
	}
	for _, session := range sessions {
		req := typespec.EvidenceUseReq{
			CheckResultId: session.ID,
			Cmd:           "whoami",
		}
		resp := typespec.EvidenceUseRes{}
		err = vulApp.VulEvidenceUse(ctx, &req, &resp)
		if err != nil {
			log.Errorf("RemoteSessionStatusCheck Error: remote session status check error: %v", err)
			session.Status = enums.SessionStatusFail
			_ = session.UpdateRemoteSession(ctx)
			continue
		}
		if resp.Result == "" {
			session.Status = enums.SessionStatusFail
			err = session.UpdateRemoteSession(ctx)
			if err != nil {
				log.Errorf("RemoteSessionStatusCheck Error: update remote session status error: %v", err)
				return
			}
		}
	}
}
