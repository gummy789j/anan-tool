package constant

import "time"

const SourcePath = "source"

type FileType string

const (
	FileTypeCSV FileType = "csv"
)

func (f FileType) String() string {
	return string(f)
}

type FileName string

const (
	MemberFile   FileName = "member"
	OriginalFile FileName = "original"
	Phase1File   FileName = "phase1"
	Phase2File   FileName = "phase2"
	WaitingFile  FileName = "waiting"
)

func (f FileName) String() string {
	return string(f)
}

type VerificationType string

const (
	VerificationTypeUnkonwn   VerificationType = "unkonwn"
	VerificationTypeThree     VerificationType = "三歲"
	VerificationTypeThreeHalf VerificationType = "三歲半"
	VerificationTypeFour      VerificationType = "四歲"
	VerificationTypeFive      VerificationType = "五歲"
)

func (a VerificationType) String() string {
	return string(a)
}

var MonthDays = func(m time.Month) int {
	switch m {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.February:
		return 28
	default:
		return 30
	}
}
