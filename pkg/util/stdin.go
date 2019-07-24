package util

import "os"

func PassedStdin() (bool, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	if (info.Mode() & os.ModeCharDevice) == 0 {
		return true, nil
	}

	return false, nil
}
