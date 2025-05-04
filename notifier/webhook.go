package notifier

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Webhook struct {
	Base

	Service string

	method          string
	contentType     string
	buildBody       func(title, message string) ([]byte, error)
	buildWebhookURL func(url string) (string, error)
	checkResult     func(status int, responseBody []byte) error
	buildHeaders    func() map[string]string
}

type webhookPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewWebhook(base *Base) *Webhook {
	base.viper.SetDefault("method", "POST")

	return &Webhook{
		Base:        *base,
		Service:     "Webhook",
		method:      base.viper.GetString("method"),
		contentType: "application/json",
		buildBody: func(title, message string) ([]byte, error) {
			return json.Marshal(webhookPayload{
				Title:   title,
				Message: message,
			})
		},
		buildHeaders: func() map[string]string {
			headers := make(map[string]string)
			for key, value := range base.viper.GetStringMapString("headers") {
				headers[key] = value
			}

			return headers
		},
		checkResult: func(status int, responseBody []byte) error {
			if status == 200 {
				return nil
			}

			return fmt.Errorf("status: %d, body: %s", status, string(responseBody))
		},
	}
}

func (s *Webhook) webhookURL() (string, error) {
	url := s.viper.GetString("url")

	if s.buildWebhookURL == nil {
		return url, nil
	}

	return s.buildWebhookURL(url)
}

func (s *Webhook) notify(title string, message string) error {
	url, err := s.webhookURL()
	if err != nil {
		return err
	}

	payload, err := s.buildBody(title, message)
	if err != nil {
		return err
	}

	slog.Info("Sending webhook notification", 
		"component", "notifier.webhook",
		"service", s.Service,
		"url", url,
		"method", s.method)
	
	req, err := http.NewRequest(s.method, url, strings.NewReader(string(payload)))
	if err != nil {
		slog.Error("Failed to create HTTP request", 
			"component", "notifier.webhook",
			"service", s.Service,
			"url", url,
			"error", err)
		return err
	}

	req.Header.Set("Content-Type", s.contentType)

	if s.buildHeaders != nil {
		headers := s.buildHeaders()
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("HTTP request failed", 
			"component", "notifier.webhook",
			"service", s.Service,
			"url", url,
			"error", err)
		return err
	}
	defer resp.Body.Close()

	var body []byte
	if resp.Body != nil {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	if s.checkResult != nil {
		err = s.checkResult(resp.StatusCode, body)
		if err != nil {
			slog.Error("Webhook response validation failed", 
				"component", "notifier.webhook",
				"service", s.Service,
				"url", url,
				"status", resp.StatusCode,
				"error", err)
			return nil
		}
	} else {
		slog.Info("Webhook response received", 
			"component", "notifier.webhook",
			"service", s.Service,
			"status", resp.StatusCode,
			"body", string(body))
	}

	slog.Info("Webhook notification sent successfully", 
		"component", "notifier.webhook",
		"service", s.Service,
		"url", url)

	return nil
}
