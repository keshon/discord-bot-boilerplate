package main

import (
	"testing"

	"github.com/keshon/discord-bot-boilerplate/internal/config"
)

func TestCreateDiscordSession_ValidToken(t *testing.T) {
	token := "MTE4MjA2Mzk1Mzc4Mjc4NDIxMQ.G1FbPm.WyCOitk1BkaVC83q2xnHTda3u5USZmBVotPyxw"
	session := createDiscordSession(token)
	if session == nil {
		t.Errorf("Expected a non-nil session, got nil")
	}
}

func TestStartRestServer(t *testing.T) {
	t.Run("RestDisabled", func(t *testing.T) {
		config := &config.Config{RestEnabled: false}
		startRestServer(config, nil)
		// Add assertion for expected behavior
	})

	t.Run("RestGinReleaseEnabled", func(t *testing.T) {
		config := &config.Config{RestEnabled: true, RestGinRelease: true}
		startRestServer(config, nil)
		// Add assertion for expected behavior
	})

	t.Run("EmptyRestHostname", func(t *testing.T) {
		config := &config.Config{RestEnabled: true, RestGinRelease: false, RestHostname: ""}
		startRestServer(config, nil)
		// Add assertion for expected behavior
	})

	t.Run("NonEmptyRestHostname", func(t *testing.T) {
		config := &config.Config{RestEnabled: true, RestGinRelease: false, RestHostname: "localhost:8080"}
		startRestServer(config, nil)
		// Add assertion for expected behavior
	})
}

func TestWaitForExitSignal(t *testing.T) {
	// Test the function waitForExitSignal
	go waitForExitSignal()
}
