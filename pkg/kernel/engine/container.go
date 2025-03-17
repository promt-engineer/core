package engine

import "sync"

var (
	bootstrap     *Bootstrap
	bootstrapOnce sync.Once
)

func PutInContainer(boot *Bootstrap) {
	bootstrapOnce.Do(func() {
		bootstrap = boot
		InitSerializer(boot.SpinFactory)
	})
}

func GetFromContainer() *Bootstrap {
	if bootstrap == nil {
		panic("container is empty")
	}

	return bootstrap
}
