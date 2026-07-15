package users_redis_cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
)

func paramsKey(filter domain.UsersFilter) (string, error) {
	raw, err := json.Marshal(filter)
	if err != nil {
		return "", fmt.Errorf("marshal users filter: %w", err)
	}

	sum := sha256.Sum256(raw)
	hashStr := hex.EncodeToString(sum[:])

	return "users:filter:by-params:" + hashStr, nil
}

func resultKey(filterID string) string {
	return "users:filter:result:" + filterID
}
