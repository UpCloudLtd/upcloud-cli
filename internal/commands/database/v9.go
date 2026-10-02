package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	upcloudv8 "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	upcloudv9 "github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
	"github.com/google/uuid"
)

const (
	databaseStateRunning = "running"
	databaseStateError   = "error"
)

// legacyDatabaseOutput preserves the database JSON/YAML contract while CRUD API calls use v9.
func legacyDatabaseOutput(db *upcloudv9.DatabaseServiceInformationResponse) (*upcloudv8.ManagedDatabase, error) {
	body, err := json.Marshal(db)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal v9 database response: %w", err)
	}
	var normalized map[string]any
	if err := json.Unmarshal(body, &normalized); err != nil {
		return nil, fmt.Errorf("cannot normalize v9 database response: %w", err)
	}
	if params, ok := normalized["service_uri_params"].(map[string]any); ok {
		for key, value := range params {
			if value != nil {
				params[key] = fmt.Sprint(value)
			}
		}
	}
	body, err = json.Marshal(normalized)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal normalized database response: %w", err)
	}

	legacy := &upcloudv8.ManagedDatabase{}
	if err := json.Unmarshal(body, legacy); err != nil {
		return nil, fmt.Errorf("cannot convert v9 database response to legacy output: %w", err)
	}
	return legacy, nil
}

func legacyDatabaseListOutput(databases []upcloudv9.DatabaseServiceInformationResponse) ([]upcloudv8.ManagedDatabase, error) {
	legacy := make([]upcloudv8.ManagedDatabase, 0, len(databases))
	for i := range databases {
		db, err := legacyDatabaseOutput(&databases[i])
		if err != nil {
			return nil, err
		}
		legacy = append(legacy, *db)
	}
	return legacy, nil
}

// databasePollInterval matches the v8 SDK wait helpers.
var databasePollInterval = 5 * time.Second

func v9Client(exec commands.Executor) (*upcloudv9.ClientWithResponses, error) {
	client := exec.V9()
	if client == nil {
		return nil, errors.New("UpCloud API v9 client is not available, run with --debug for details")
	}
	return client, nil
}

func deref[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

func parseDatabaseUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return id, fmt.Errorf("invalid database UUID %q: %w", value, err)
	}
	return id, nil
}

func uuidString(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

func mapString(m *map[string]any, key string) string {
	if m == nil {
		return ""
	}
	switch v := (*m)[key].(type) {
	case nil:
		return ""
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

type propertySchema struct {
	// Type is a JSON schema type: a string, or a list such as ["string", "null"].
	Type any `json:"type"`
}

func (p propertySchema) allowsString() bool {
	switch t := p.Type.(type) {
	case string:
		return t == "string"
	case []any:
		return slices.Contains(t, any("string"))
	}
	return false
}

// getServiceTypeProperties decodes the response itself, as the generated v9 model cannot decode list-valued property types.
func getServiceTypeProperties(ctx context.Context, client *upcloudv9.ClientWithResponses, dbType string) (map[string]propertySchema, error) {
	res, err := client.GetDatabaseType(ctx, upcloudv9.GetDatabaseTypeServiceTypeName(dbType))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, apierror.FromResponse(res.StatusCode, body)
	}

	var serviceType struct {
		Properties map[string]propertySchema `json:"properties"`
	}
	if err := json.Unmarshal(body, &serviceType); err != nil {
		return nil, fmt.Errorf("cannot parse database type %s: %w", dbType, err)
	}
	return serviceType.Properties, nil
}

// waitForDatabaseState polls like the v8 SDK, retrying up to three consecutive server errors.
func waitForDatabaseState(ctx context.Context, client *upcloudv9.ClientWithResponses, id uuid.UUID, state string) error {
	ticker := time.NewTicker(databasePollInterval)
	defer ticker.Stop()

	serverErrors := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			res, err := client.GetDatabaseWithResponse(ctx, id)
			if err != nil {
				return err
			}
			if res.JSON200 == nil {
				if res.StatusCode() >= http.StatusInternalServerError && serverErrors < 3 {
					serverErrors++
					continue
				}
				return apierror.FromResponse(res.StatusCode(), res.Body)
			}
			serverErrors = 0

			switch deref(res.JSON200.State) {
			case state:
				return nil
			case databaseStateError:
				return stateError(deref(res.JSON200.StateError))
			}
		}
	}
}

func stateError(stateErrors map[string]string) error {
	details := make([]string, 0, len(stateErrors))
	for _, key := range slices.Sorted(maps.Keys(stateErrors)) {
		details = append(details, fmt.Sprintf("%s (%s)", stateErrors[key], key))
	}
	return fmt.Errorf("managed database entered error state: %s", strings.Join(details, ", "))
}
