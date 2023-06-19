package parser

type PlaintextParser struct {
	startIdx int
	endIdx   int
}

func NewPlaintextParser(startIdx int, endIdx int) *PlaintextParser {
	return &PlaintextParser{
		startIdx: startIdx,
		endIdx:   endIdx,
	}
}

func (ptp *PlaintextParser) Parse(data []byte) (string, error) {
	text := string(data)

	if ptp.startIdx >= 0 && ptp.startIdx <= len(text) {
		if ptp.endIdx >= 0 && ptp.endIdx <= len(text) {
			return text[ptp.startIdx:ptp.endIdx], nil
		}
		return text[ptp.startIdx:], nil

	}

	if ptp.endIdx >= 0 && ptp.endIdx <= len(text) {
		return text[:ptp.endIdx], nil
	}
	return text, nil
}
