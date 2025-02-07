package writer

type FileExt uint
type ExtBitset uint8

const (
	None     FileExt = 0
	Standard FileExt = 1
	TXT      FileExt = 2
	CSV      FileExt = 3
)

func extFromString(extStr string) FileExt {
	switch extStr {
	case "std":
		return Standard
	case "txt":
		return TXT
	case "csv":
		return CSV
	default:
		return None
	}
}
