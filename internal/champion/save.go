package champion

import (
	"fmt"

	"github.com/5pots-com/cli/internal/common"
)

func (c Champion) SaveToFile(dir, fileName string, data []byte) error {
	if err := common.SaveToFile(dir, fileName, data); err != nil {
		return fmt.Errorf("error while saving champion to file: %v", err)
	}

	return nil
}
