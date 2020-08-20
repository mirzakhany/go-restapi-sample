package registry

import (
	"github.com/mirzakhany/rest_api_sample/pkg/projectx"
	"sort"
)

var items registryItems

type registryItems []registryItem

type registryItem struct {
	function    func(ctx *projectx.Ctx) error
	priority    int
	failOnError bool
}

func (r registryItems) Len() int {
	return len(r)
}

func (r registryItems) Less(i, j int) bool {
	return r[i].priority < r[j].priority
}

func (r registryItems) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func Register(function func(ctx *projectx.Ctx) error, priority int, failOnError bool) {
	items = append(items, registryItem{
		function:    function,
		priority:    priority,
		failOnError: failOnError,
	})
}

func Flush() {
	items = nil
}

func Run(ctx *projectx.Ctx) (int, []error) {

	// sort items
	sort.Sort(items)

	var errors []error
	var err error

	// iterate over items and run each item.
	for i, item := range items {
		err = item.function(ctx)
		if err != nil {
			errors = append(errors, err)
		}
		if item.failOnError {
			return i + 1, errors
		}
	}
	return len(items), errors
}
