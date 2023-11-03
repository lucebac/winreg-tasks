//go:build !windows

package providers

import "fmt"

func GetNativeSystemProvider() (DataProvider, error) {
	return nil, fmt.Errorf("no native data provider available on this platform")
}
