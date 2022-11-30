package brest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

// Engine structure
type Engine struct {
	config *Config
}

// NewEngine constructs Engine
func NewEngine(config *Config) *Engine {
	e := new(Engine)
	e.config = config
	e.config.InfoLogger().Printf("Brest configuration: %v\n", e.config)
	return e
}

// Config gets config
func (e *Engine) Config() *Config {
	return e.config
}

// Execute executes a rest query
func (e *Engine) Execute(restQuery *RestQuery) (interface{}, error) {
	if restQuery.Debug {
		e.Config().InfoLogger().Printf("Execution request: %v\n", restQuery)
	}
	resource, err := e.getResource(restQuery)
	if err != nil {
		return nil, NewErrorFromCause(err)
	}
	if resource.Action()&restQuery.Action == 0 {
		return nil, NewErrorForbbiden(fmt.Sprintf("query %v not authorized for resource %v", restQuery, resource))
	}
	var entity interface{}
	var elem reflect.Value
	if restQuery.Action == Get {
		if restQuery.Key != "" {
			elem = reflect.New(resource.ResourceType()).Elem()
			entity = elem.Addr().Interface()
			if err = setPk(e.Config().DB(), resource.ResourceType(), elem, restQuery.Key); err != nil {
				return nil, NewErrorFromCause(err)
			}
		} else {
			sliceType := reflect.MakeSlice(reflect.SliceOf(resource.ResourceType()), 0, 0).Type()
			entity = reflect.New(sliceType).Interface()
		}
	} else if restQuery.Action == Post {
		if restQuery.Key != "" {
			return nil, NewErrorBadRequest("action 'Post': key is forbidden")
		}
		elem = reflect.New(resource.ResourceType()).Elem()
		entity = elem.Addr().Interface()
		if err = e.Deserialize(restQuery, resource, entity); err != nil {
			return nil, NewErrorFromCause(err)
		}
	} else if restQuery.Action == Put {
		if restQuery.Key == "" {
			return nil, NewErrorBadRequest("action 'Put': key is mandatory")
		}
		elem = reflect.New(resource.ResourceType()).Elem()
		entity = elem.Addr().Interface()
		if err = e.Deserialize(restQuery, resource, entity); err != nil {
			return nil, NewErrorFromCause(err)
		}
		setPk(e.Config().DB(), resource.ResourceType(), elem, restQuery.Key)
	} else if restQuery.Action == Patch {
		if restQuery.Key == "" {
			return nil, NewErrorBadRequest("action 'Patch': key is mandatory")
		}
		elem = reflect.New(resource.ResourceType()).Elem()
		entity = elem.Addr().Interface()
		if err = setPk(e.Config().DB(), resource.ResourceType(), elem, restQuery.Key); err != nil {
			return nil, NewErrorFromCause(err)
		}
	} else if restQuery.Action == Delete {
		if restQuery.Key == "" {
			return nil, NewErrorBadRequest("action 'Delete': key is mandatory")
		}
		elem = reflect.New(resource.ResourceType()).Elem()
		entity = elem.Addr().Interface()
		if err = setPk(e.Config().DB(), resource.ResourceType(), elem, restQuery.Key); err != nil {
			return nil, NewErrorFromCause(err)
		}
	} else {
		return nil, NewErrorBadRequest(fmt.Sprintf("unknow action '%v'", restQuery.Action))
	}

	var ctx context.Context
	if restQuery.Request != nil {
		ctx = restQuery.Request.Context()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if resource.beforeHook != nil {
		if err = resource.beforeHook(ctx, resource.name, resource.Action(), entity); err != nil {
			return nil, NewErrorFromCause(err)
		}
	}

	ctx = ContextWithDb(ctx, e.Config().DB())

	executor := NewExecutor(e.Config(), restQuery, entity, restQuery.Debug)

	if restQuery.Action == Get {
		if restQuery.Key != "" {
			err = executor.Execute(ctx, executor.GetOneExecFunc())
		} else {
			err = executor.Execute(ctx, executor.GetSliceExecFunc())
		}
	} else if restQuery.Action == Post {
		err = executor.Execute(ctx, executor.InsertExecFunc())
	} else if restQuery.Action == Put {
		err = executor.Execute(ctx, executor.UpdateExecFunc())
	} else if restQuery.Action == Patch {
		err = executor.Execute(ctx, executor.GetOneExecFunc())
		if err == nil {
			err = e.Deserialize(restQuery, resource, entity)
		}
		if err == nil {
			err = setPk(e.Config().DB(), resource.ResourceType(), elem, restQuery.Key)
		}
		if err == nil {
			err = executor.Execute(ctx, executor.UpdateExecFunc())
		}
	} else if restQuery.Action == Delete {
		err = executor.Execute(ctx, executor.DeleteExecFunc())
	}
	if err != nil {
		return nil, NewErrorFromCause(err)
	}
	if resource.afterHook != nil {
		if err = resource.afterHook(ctx, resource.name, resource.Action(), entity); err != nil {
			return nil, NewErrorFromCause(err)
		}
	}
	if restQuery.Debug {
		e.Config().InfoLogger().Printf("Execution result: %v\n", entity)
	}
	if restQuery.Action == Get && restQuery.Key == "" {
		return NewPage(executor.entity, executor.count, restQuery), nil
	}
	return executor.entity, nil
}

// Deserialize deserializes data into entity
func (e *Engine) Deserialize(restQuery *RestQuery, resource *Resource, entity interface{}) error {
	if regexp.MustCompile("[+-/]json($|[+-;])").MatchString(restQuery.ContentType) {
		if err := json.Unmarshal(restQuery.Content, entity); err != nil {
			return NewErrorFromCause(err)
		}
	} else if regexp.MustCompile("[+-/]form($|[+-;])").MatchString(restQuery.ContentType) {
		table := e.config.db.Table(resource.ResourceType())
		keyValues := strings.Split(string(restQuery.Content), "&")
		elem := reflect.ValueOf(entity).Elem()
		for _, keyValue := range keyValues {
			parts := strings.Split(keyValue, "=")
			if parts != nil && len(parts) == 2 {
				found := false
				for _, field := range table.Fields {
					if field.GoName == parts[0] {
						field.ScanValue(elem, parts[1])
						found = true
					}
				}
				if !found {
					for _, field := range table.Fields {
						if field.Name == parts[0] {
							field.ScanValue(elem, parts[1])
							found = true
						}
					}
				}
			}
		}
	} else if regexp.MustCompile("[+-/](msgpack|messagepack)($|[+-])").MatchString(restQuery.ContentType) {
		decoder := msgpack.NewDecoder(bytes.NewReader(restQuery.Content))
		decoder.SetCustomStructTag("json")
		if err := decoder.Decode(entity); err != nil {
			return *NewErrorFromCause(err)
		}
	} else {
		return NewErrorBadRequest(fmt.Sprintf("Unknown content type '%v'", restQuery.ContentType))
	}
	if restQuery.Debug {
		e.Config().InfoLogger().Printf("Serialized response in %v: %v\n", restQuery.ContentType, entity)
	}
	return nil
}

func (e *Engine) getResource(restQuery *RestQuery) (*Resource, error) {
	if restQuery.Resource == "" {
		return nil, NewErrorBadRequest("resource is mandatory")
	}
	resource := e.config.GetResource(restQuery.Resource)
	if resource == nil {
		e.Config().ErrorLogger().Printf("Resource '%v' not defined in engine configuration", restQuery.Resource)
		e.Config().ErrorLogger().Printf("Configuration: '%v'", e.config)
		return nil, NewErrorBadRequest(fmt.Sprintf("resource '%v' not defined in engine configuration", restQuery.Resource))
	}
	return resource, nil
}
