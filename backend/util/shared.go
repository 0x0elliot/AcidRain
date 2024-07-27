package util 

import (
	"encoding/base64"
)

func Contains(arr []string, str string) bool {
	
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func ContainsInt64(arr []int64, num int64) bool {
	for _, a := range arr {
		if a == num {
			return true
		}
	}
	return false
}

func IsBase64Image(base64Image string) bool {
	return "data:image/" == base64Image[:11]
}

// CalculateBase64ImageSizeMB takes a base64 encoded string and returns its size in megabytes
func CalculateBase64ImageSizeMB(base64String string) (float64, error) {
	// Decode the base64 string
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return 0, err
	}

	// Calculate the size in bytes
	sizeInBytes := len(data)

	// Convert the size to megabytes (1 MB = 1024 * 1024 bytes)
	sizeInMB := float64(sizeInBytes) / (1024 * 1024)

	return sizeInMB, nil
}

