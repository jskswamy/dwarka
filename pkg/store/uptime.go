package store

import (
	"encoding/json"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/gateway"
	"path"
	"time"
)

const (
	uptimePath = "status/server"
	startTime  = "startTime"
)

// Uptime returns the uptime information from store
func (ps PersistentStore) Uptime() (gateway.Status, error) {
	data, err := ps.get(ps.uptimeRootPath(), gateway.Status{startTime: now()})
	if err != nil {
		return nil, err
	}

	result := gateway.Status{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RefreshUptime updates the uptime in store
func (ps PersistentStore) RefreshUptime() error {
	status := gateway.Status{startTime: now()}
	return ps.putJSON(ps.uptimeRootPath(), status)
}

func (ps PersistentStore) uptimeRootPath() string {
	return path.Join(ps.path, uptimePath)
}

func now() string {
	return time.Now().Local().Format(time.RFC822)
}
