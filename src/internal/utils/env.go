package utils

import (
	"os"
	"pangya/src/internal/logger"
	"strconv"
)

func GetStringEnv(name string) string {
	return os.Getenv(name)
}

func GetIntEnv(name string) int {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	return val
}

func GetUint16Env(name string) uint16 {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	return uint16(val)
}

func GetUint32Env(name string) uint32 {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	return uint32(val)
}
