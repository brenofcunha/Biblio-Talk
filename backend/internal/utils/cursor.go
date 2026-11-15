package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type FeedCursor struct {
	CreatedAt time.Time
	ID        string
}

func EncodeFeedCursor(createdAt time.Time, id string) string {
	raw := fmt.Sprintf("%d|%s", createdAt.UnixMilli(), id)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func DecodeFeedCursor(cursor string) (FeedCursor, error) {
	var c FeedCursor
	decoded, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return c, fmt.Errorf("cursor inválido")
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return c, fmt.Errorf("cursor inválido")
	}
	var millis int64
	_, err = fmt.Sscanf(parts[0], "%d", &millis)
	if err != nil {
		return c, fmt.Errorf("cursor inválido")
	}
	c.CreatedAt = time.UnixMilli(millis)
	c.ID = parts[1]
	return c, nil
}
