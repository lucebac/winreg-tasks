package providers

func GetNativeSystemProvider() (DataProvider, error) {
	return NewWindowsRegistryProvider()
}
