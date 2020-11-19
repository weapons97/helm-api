package utils

import "k8s.io/apimachinery/pkg/util/yaml"

func YamlToJson(ym []byte) (json []byte, err error) {
	return yaml.ToJSON(ym)
}
