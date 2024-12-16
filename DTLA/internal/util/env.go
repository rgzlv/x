package util

import (
	"os"
	"strconv"
)

func EnvGetInt(name string) (int, error) {
	envStr, ok := os.LookupEnv(name)
	if !ok {
		envStr = "0"
	}

	return strconv.Atoi(envStr)
}
