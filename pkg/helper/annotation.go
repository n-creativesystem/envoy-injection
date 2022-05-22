package helper

import "strings"

type Annotation interface {
	GetKey(name ...string) string
	GetValue(names ...string) string
	GetValueOrDefault(name, def string) string
	IsEnabled() bool
}

type annotationImpl struct {
	prefix      string
	annotations map[string]string
}

// NewAnnotationHelper keys = ncs.extensions.keys...k8s.io
func NewAnnotationHelper(annotation map[string]string, keys ...string) Annotation {
	return annotationImpl{
		prefix:      annotationName(keys...),
		annotations: annotation,
	}
}

func (a annotationImpl) GetKey(names ...string) string {
	n := append([]string{a.prefix}, names...)
	return strings.Join(n, "/")
}

func (a annotationImpl) GetValue(names ...string) string {
	value, ok := a.annotations[a.GetKey(names...)]
	if !ok {
		return ""
	}
	return value
}

func (a annotationImpl) GetValueOrDefault(name, def string) string {
	v := a.GetValue(name)
	if v == "" {
		return def
	}
	return v
}

func (a annotationImpl) IsEnabled() bool {
	return a.GetValue("enabled") == "true"
}
