package linenotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type notifyOption func(url.Values)

func WithSticker(stickerPackageID, stickerID int) notifyOption {
	return func(v url.Values) {
		v.Set("stickerPackageId", strconv.Itoa(stickerPackageID))
		v.Set("stickerId", strconv.Itoa(stickerID))
	}
}

func Silent() notifyOption {
	return func(v url.Values) {
		v.Set("notificationDisabled", "true")
	}
}

func Notify(ctx context.Context, token string, message string, options ...notifyOption) error {
	form := url.Values{}
	for _, o := range options {
		o(form)
	}
	form.Set("message", message)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://notify-api.line.me/api/notify", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var nr notifyResponse
	if err := json.NewDecoder(req.Body).Decode(&nr); err != nil {
		return err
	}

	if nr.Status != 200 {
		return fmt.Errorf("status=%d, message=%s", nr.Status, nr.Message)
	}

	return nil
}

type notifyResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
