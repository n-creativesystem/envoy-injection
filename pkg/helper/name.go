package helper

import "strings"

const (
	ComponentName = "ncs-kubernetes-extensions"
)

func annotationName(keys ...string) string {
	key := strings.Join(keys, ".")
	return strings.Join([]string{"ncs.extensions", key, "k8s.io"}, ".")
}
