package macos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

func Notify(text string) error {
	j, err := json.Marshal(text)
	if err != nil {
		return err
	}
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification %s with title "gh status" sound name "Boop"`, j))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, stderr.String())
	}
	return nil
}
