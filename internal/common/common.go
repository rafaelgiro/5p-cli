package common

type Patch string

const (
	PBE    Patch = "pbe"
	Latest Patch = "latest"
)

func Validate(patch Patch) bool {
	switch patch {
	case PBE, Latest:
		return true
	}
	return false
}
