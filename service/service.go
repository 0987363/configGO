package service

import "os"

type Service struct {
	CloseHook []func()
}

var service *Service

func init() {
	service = &Service{}
}

func AddCloseHook(f func()) {
	service.CloseHook = append(service.CloseHook, f)
}

func Exit(code int) {
	for _, f := range service.CloseHook {
		f()
	}
	os.Exit(code)
}
