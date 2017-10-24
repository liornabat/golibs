package logging

import (
	"github.com/stretchr/testify/require"
	"testing"

	"errors"
	"fmt"
	"os"
)

func TestMain(m *testing.M) {
	InitLoggers()
	os.Exit(m.Run())
}

func TestL(t *testing.T) {
	require := require.New(t)
	log := NewLogger("test")
	require.NotNil(log.log)
	log.SetDebug(true)
	log.Error(errors.New("some error"), "some text error 1")
	log.Error(errors.New("some error"), "some text error 1", "some text error 2")
	log.Error(nil, "some text error 1", "some text error 2")
	log.Error(nil, "some text error 1")
	log.Error(nil)

	log.Info("some text info ")
	log.Debug("some text debug ", nil)

	log.Panic(errors.New("some error"), "sad")

}

func TestHook(t *testing.T) {
	require := require.New(t)
	log := NewLogger("test")
	require.NotNil(log.log)
	log.SetHook(func(msg string) {
		fmt.Println("testing hook ", msg)
	})
	log.Info("testing")
}
