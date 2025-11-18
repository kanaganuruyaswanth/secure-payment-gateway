package main

import (
	"encoding/json"
	"time"
)

func CreateAuditLog(entityType, entityID, action string, metadata interface{}) {
	data, _ := json.Marshal(metadata)

	log := AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Metadata:   string(data),
		CreatedAt:  time.Now(),
	}

	DB.Create(&log)
}
