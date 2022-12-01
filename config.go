package brest

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/uptrace/bun"
)

// BeforeHook defines before execution callback function
type BeforeHook func(ctx context.Context, restQuery *RestQuery, entity interface{}) error

// AfterHook defines after execution callback function
type AfterHook func(ctx context.Context, restQuery *RestQuery, entity interface{}) error

// Resource structure
type Resource struct {
	name         string
	resourceType reflect.Type
	action       Action
	beforeHook   BeforeHook
	afterHook    AfterHook
}

func (r *Resource) String() string {
	var str string
	str = fmt.Sprintf("<name=%v resourceType=%v action=%v>", r.name, r.resourceType, r.action)
	return str
}

// Name returns name
func (r *Resource) Name() string {
	return r.name
}

// ResourceType returns resourceType
func (r *Resource) ResourceType() reflect.Type {
	return r.resourceType
}

// Action returns action
func (r *Resource) Action() Action {
	return r.action
}

// NewResource constructs Resource
func NewResource(name string, entity interface{}, action Action) *Resource {
	r := new(Resource)
	r.name = name
	r.resourceType = reflect.TypeOf(entity)
	if r.resourceType.Kind() == reflect.Ptr {
		r.resourceType = r.resourceType.Elem()
	}
	r.action = action
	return r
}

// NewResourceWithHooks constructs Resource with hooks
func NewResourceWithHooks(name string, entity interface{}, action Action, beforeHook BeforeHook, afterHook AfterHook) *Resource {
	r := new(Resource)
	r.name = name
	r.resourceType = reflect.TypeOf(entity)
	if r.resourceType.Kind() == reflect.Ptr {
		r.resourceType = r.resourceType.Elem()
	}
	r.action = action
	r.beforeHook = beforeHook
	r.afterHook = afterHook
	return r
}

// Config structure
type Config struct {
	prefix             string
	db                 *bun.DB
	resources          map[string]*Resource
	defaultContentType string
	defaultAccept      string
	infoLogger         *log.Logger
	errorLogger        *log.Logger
}

func (c *Config) String() string {
	return fmt.Sprintf("version=%v db=%v resources=%v", Version(), c.db, c.resources)
}

// AddResource adds resource
func (c *Config) AddResource(resource *Resource) {
	elem := reflect.New(resource.ResourceType()).Elem()
	entity := elem.Addr().Interface()
	c.db.RegisterModel(entity)
	c.resources[resource.Name()] = resource
}

// GetResource gets resource
func (c *Config) GetResource(resourceName string) *Resource {
	return c.resources[resourceName]
}

// SetPrefix sets prefix
func (c *Config) SetPrefix(prefix string) {
	c.prefix = prefix
	if c.prefix == "" {
		c.prefix = "/"
	}
	if !strings.HasPrefix(c.prefix, "/") {
		c.prefix = "/" + c.prefix
	}
	if !strings.HasSuffix(c.prefix, "/") {
		c.prefix = c.prefix + "/"
	}
}

// Prefix gets prefix
func (c *Config) Prefix() string {
	return c.prefix
}

// SetDefaultContentType sets defaultContentType
func (c *Config) SetDefaultContentType(defaultContentType string) {
	c.defaultContentType = defaultContentType
}

// DefaultContentType gets defaultContentType
func (c *Config) DefaultContentType() string {
	return c.defaultContentType
}

// SetDefaultAccept sets defaultAccept
func (c *Config) SetDefaultAccept(defaultAccept string) {
	c.defaultAccept = defaultAccept
}

// DefaultAccept gets defaultAccept
func (c *Config) DefaultAccept() string {
	return c.defaultAccept
}

// DB gets db
func (c *Config) DB() *bun.DB {
	return c.db
}

// SetInfoLogger sets info logger
func (c *Config) SetInfoLogger(logger *log.Logger) {
	c.infoLogger = logger
}

// InfoLogger gets info logger
func (c *Config) InfoLogger() *log.Logger {
	return c.infoLogger
}

// SetErrorLogger sets error logger
func (c *Config) SetErrorLogger(logger *log.Logger) {
	c.errorLogger = logger
}

// ErrorLogger gets error logger
func (c *Config) ErrorLogger() *log.Logger {
	return c.errorLogger
}

// NewConfig constructs Config
func NewConfig(prefix string, db *bun.DB) *Config {
	c := new(Config)
	c.SetPrefix(prefix)
	c.db = db
	c.resources = make(map[string]*Resource)
	c.defaultContentType = Json
	c.defaultAccept = Json
	c.infoLogger = log.New(os.Stdout, " INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	c.errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return c
}
