package utils

import "golang.org/x/net/html"

func FindAttribute(attributes []html.Attribute, namespace, key string) (found bool, value string) {
	for _, attribute := range attributes {
		if attribute.Namespace == namespace &&
			attribute.Key == key {
			return true, attribute.Val
		}
	}
	return false, ""
}

func HasAttribute(attributes []html.Attribute, namespace, key, value string) bool {
	for _, attribute := range attributes {
		if attribute.Namespace == namespace &&
			attribute.Key == key &&
			attribute.Val == value {
			return true
		}
	}
	return false
}
