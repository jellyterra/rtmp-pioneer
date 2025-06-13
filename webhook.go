package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WebhookConfig struct {
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	Headers      map[string]string `json:"headers"`
	Cookies      map[string]string `json:"cookies"`
	Body         string            `json:"body"`
	BodyEncoding string            `json:"body_encoding"`
}

func DoWebhook(h *Handler, config *WebhookConfig) error {

	var requestBody io.Reader

	switch config.BodyEncoding {
	case "string":
		requestBody = strings.NewReader(config.Body)
	case "base64":
		b, err := base64.StdEncoding.DecodeString(config.Body)
		if err != nil {
			return err
		}
		requestBody = bytes.NewReader(b)
	default:
		return fmt.Errorf("unexpected body encoding: %s", config.BodyEncoding)
	}

	req, err := http.NewRequest(config.Method, config.URL, requestBody)
	if err != nil {
		return err
	}

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	for key, value := range config.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  key,
			Value: value,
		})
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	b, _ := io.ReadAll(resp.Body)
	h.Logln("Webhook response", resp.StatusCode, "text:", string(b))
	_ = resp.Body.Close()

	return nil
}
