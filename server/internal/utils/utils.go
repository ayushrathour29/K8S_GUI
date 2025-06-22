package utils

import "k8s.io/apimachinery/pkg/util/intstr"

// Helper function to convert int32 to IntOrString
func IntstrFromInt(i int32) intstr.IntOrString {
	return intstr.FromInt(int(i))
}
