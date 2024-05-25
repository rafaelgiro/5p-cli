package champion

import (
	"fmt"
	"os"
)

func (c *Champion) Load(dir string, p string) ([]byte, error) {
	fileName := fmt.Sprintf("%s/%s/%s", dir, c.Name, p)

	f, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", fileName, err)
	}

	return f, nil
}
