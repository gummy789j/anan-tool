package model

import (
	"time"

	"github.com/gummy789j/anan-tool/constant"
)

type Phase1 struct {
	Name  string `csv:"幼兒姓名"`
	Birth string `csv:"出生年月日"`
}

type Phase2 struct {
	Name  string `csv:"幼兒姓名"`
	Birth string `csv:"出生年月日"`
}

type Waiting struct {
	Name  string `csv:"姓名"`
	Birth string `csv:"生日"`
}

type Original struct {
	Name  string `csv:"幼兒姓名"`
	Birth string `csv:"出生年月日"`
}

type Member struct {
	Name string `csv:"姓名"`
}

type MemberInfo struct {
	Name             string
	Birth            time.Time
	Age              []int
	VerificationType constant.VerificationType
}

/*
編號, 姓名, 生日(民國年/月/日), 歲數(x年x月x日), 檢核檔(age type)
*/
type Result struct {
	SerialNum        string `csv:"編號"`
	Name             string `csv:"姓名"`
	Birth            string `csv:"生日"`
	Age              string `csv:"年紀"`
	VerificationType string `csv:"檢核檔"`
}

type CountTypeResult struct {
	VerificationTypeThree     int `csv:"三歲"`
	VerificationTypeThreeHalf int `csv:"三歲半"`
	VerificationTypeFour      int `csv:"四歲"`
	VerificationTypeFive      int `csv:"五歲"`
}
