// Package services
// @Author bcy2007  2024/12/12 16:58
package services

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
)

type FrpsService struct{}

func (srv *FrpsService) RunningFrpsService(ctx context.Context) error {
	status, _ := srv.FindRunningFrpsService(ctx)
	if status == true {
		return errors.New("frps service is running")
	}
	paramsList := []string{"-c", "./bin/frps.ini"}
	cmd := exec.CommandContext(ctx, "./bin/frps", paramsList...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Info("RunningFrpsService cmd StdoutPipe err:", err)
		return err
	}
	defer stdout.Close()
	if err = cmd.Start(); err != nil {
		log.Info("RunningFrpsService cmd Start err:", err)
		return err
	}
	return cmd.Wait()
}

func (srv *FrpsService) FindRunningFrpsService(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "pgrep", "frps")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("FindRunningFrpsService get combined output %v, error: %v", string(output), err)
	}
	log.Infof("pgrep frps: %s", string(output))
	processNumber, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return false, fmt.Errorf("FindRunningFrpsService get process number error: %v", err)
	}
	log.Infof("FindRunningFrpsService process number: %d", processNumber)
	return true, nil
}

func (srv *FrpsService) StopFrpsService(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "pkill", "frps")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
	fmt.Println("pkill frps: ", string(out))
	return nil
}
