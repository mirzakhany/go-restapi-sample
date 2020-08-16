package registry

import (
	"context"
	"fmt"
	"testing"
)

func TestFlush(t *testing.T) {

	ctx := context.Background()
	// function without error
	a := func(ctx context.Context) error {
		return fmt.Errorf("func error")
	}

	b := func(ctx context.Context) error {
		return fmt.Errorf("func error")
	}

	Register(a, 0, false)
	Register(b, 0, false)

	Flush()

	runs, errors := Run(ctx)
	if len(errors) > 0 {
		t.Errorf("expected no error returned: %d", len(errors))
	}

	if runs != 0 {
		t.Errorf("expected to run no task but runs: %d", runs)
	}
}

func TestRegister(t *testing.T) {

	ctx := context.Background()

	// function without error
	a := func(ctx context.Context) error {
		return nil
	}

	b := func(ctx context.Context) error {
		return fmt.Errorf("func error")
	}

	Register(a, 0, false)
	Register(b, 0, false)

	runs, errors := Run(ctx)
	if len(errors) > 1 {
		t.Errorf("expected one error returned: %d", len(errors))
	}

	if runs != 2 {
		t.Errorf("expected to run two task but runs: %d", runs)
	}

	c := func(ctx context.Context) error {
		return fmt.Errorf("func c error")
	}

	Flush()

	Register(a, 0, false)
	Register(b, 0, false)
	Register(c, 0, false)

	runs, errors = Run(ctx)
	if len(errors) != 2 {
		t.Errorf("expected one error returned: %d", len(errors))
	}

	if runs != 3 {
		t.Errorf("expected to run three task but runs: %d", runs)
	}

	Flush()

	Register(a, 0, false)
	Register(b, 0, true)
	Register(c, 0, false)

	runs, errors = Run(ctx)
	if len(errors) != 1 {
		t.Errorf("expected one error returned: %d", len(errors))
	}

	if runs != 2 {
		t.Errorf("expected to run two task but runs: %d", runs)
	}
}
