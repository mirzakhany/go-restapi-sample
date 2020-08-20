package projectx

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {

	parent := context.Background()

	pctx := New(parent)

	if pctx == nil {
		t.Error("project context is nil")
	}
}

func TestCtx_Set(t *testing.T) {
	parent := context.Background()
	pctx := New(parent)

	pctx.Set("test", "this is a value")
	v, exists := pctx.Get("test")
	if !exists {
		t.Error("key not found")
	}

	if v.(string) != "this is a value" {
		t.Error("value is not equal")
	}
}

func TestCtx_Get(t *testing.T) {
	parent := context.Background()
	pctx := New(parent)

	pctx.Set("test_get", "this is a value")
	v, exists := pctx.Get("test_get")
	if !exists {
		t.Error("key not found")
	}

	if v.(string) != "this is a value" {
		t.Error("value is not equal")
	}
}

func TestCtx_Get_NotExist(t *testing.T) {
	parent := context.Background()
	pctx := New(parent)

	v, exists := pctx.Get("test_not_exist")
	if exists {
		t.Error("exist should be false but is true")
	}

	if v!=nil {
		t.Errorf("value should be nil but is %v", v)
	}
}