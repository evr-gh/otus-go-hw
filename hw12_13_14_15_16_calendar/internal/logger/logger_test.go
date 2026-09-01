package logger

import (
	"bytes"
	"strings"
	"testing"
)

const (
	debugMsg   string = "отладочное сообщение"
	infoMsg    string = "информационное сообщение"
	warningMsg string = "предупреждение"
	errorMsg   string = "ошибка"
)

func TestLogger(t *testing.T) {
	testCases := []struct {
		level            LogLevel
		debugMessage     string
		infoMessage      string
		warningMessage   string
		errorMessage     string
		expectedMessages []string
	}{
		{
			DEBUG,
			debugMsg,
			infoMsg,
			warningMsg,
			errorMsg,
			[]string{debugMsg, infoMsg, warningMsg, errorMsg},
		},
		{
			INFO,
			debugMsg,
			infoMsg,
			warningMsg,
			errorMsg,
			[]string{infoMsg, warningMsg, errorMsg},
		},
		{
			WARNING,
			debugMsg,
			infoMsg,
			warningMsg,
			errorMsg,
			[]string{warningMsg, errorMsg},
		},
		{
			ERROR,
			debugMsg,
			infoMsg,
			"предупреждение",
			errorMsg,
			[]string{errorMsg},
		},
	}
	for _, testCase := range testCases {
		outputInto := &bytes.Buffer{}
		logger := New(testCase.level, outputInto)
		logger.Debug("%s", testCase.debugMessage)
		logger.Info("%s", testCase.infoMessage)
		logger.Warning("%s", testCase.warningMessage)
		logger.Error("%s", testCase.errorMessage)
		output := outputInto.String()
		for _, expected := range testCase.expectedMessages {
			if !strings.Contains(output, expected) {
				t.Errorf("Выввод %s (%q) не содержит %s\n", logger.level, output, expected)
			}
		}
	}
}

func TestOutputFormat(t *testing.T) {
	expected := "Тестовое сообщение с несколькими полями s=\"текст\" i=12345\n"
	outputBuff := &bytes.Buffer{}
	logger := New("INFO", outputBuff)
	logger.Info("Тестовое сообщение с несколькими полями s=%q i=%v\n", "текст", 12345)
	output := outputBuff.String()
	if !strings.HasSuffix(output, expected) {
		t.Errorf("Выввод %s (%q) не содержит %q\n", logger.level, output, expected)
	}
	if !strings.HasPrefix(output, "INFO") {
		t.Errorf("Выввод %s (%q) не содержит %q\n", logger.level, output, "INFO")
	}
}
