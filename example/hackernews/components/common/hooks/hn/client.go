package hn

import (
	"github.com/munrocape/hn/hnclient"
)

func UseHackerNews() *hnclient.Client {
	c := hnclient.NewClient()
	return c
}
