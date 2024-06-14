package log

func GetLogDriver(logDriverType string) (driver Driver, err error) {
	switch logDriverType {
	case DriverTypeFile:
		driver, err = GetFileLogDriver()
		if err != nil {
			return driver, err
		}
	case DriverTypeMongo:
		return driver, ErrNotImplemented
	case DriverTypeEs:
		return driver, ErrNotImplemented
	default:
		return driver, ErrInvalidType
	}
	return driver, nil
}
