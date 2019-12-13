// +build !linux

package api

func maxConcurrency() int {
	return 128
}
