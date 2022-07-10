package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/gummy789j/anan-tool/constant"
	"github.com/gummy789j/anan-tool/model"
	"github.com/gummy789j/anan-tool/util"
)

var location, _ = time.LoadLocation("Asia/Taipei")
var targetDate = time.Date(2022, time.September, 30, 0, 0, 0, 0, location)

var countVerificationType = map[constant.VerificationType]int{}

var memberMap = map[string]bool{}

func main() {
	memberPath := getPath(constant.MemberFile, constant.FileTypeCSV, constant.SourcePath)
	phase1Path := getPath(constant.Phase1File, constant.FileTypeCSV, constant.SourcePath)
	phase2Path := getPath(constant.Phase2File, constant.FileTypeCSV, constant.SourcePath)
	waitingPath := getPath(constant.WaitingFile, constant.FileTypeCSV, constant.SourcePath)
	originalPath := getPath(constant.OriginalFile, constant.FileTypeCSV, constant.SourcePath)

	memberInfo, err := getMemberFromCSV(memberPath)
	if err != nil {
		return
	}
	phase1Info, err := getPhase1FromCSV(phase1Path)
	if err != nil {
		return
	}
	phase2Info, err := getPhase2FromCSV(phase2Path)
	if err != nil {
		return
	}
	waitingInfo, err := getWaitingFromCSV(waitingPath)
	if err != nil {
		return
	}
	originalInfo, err := getOriginalFromCSV(originalPath)
	if err != nil {
		return
	}

	for _, info := range memberInfo {
		_, ok := memberMap[info.Name]
		if !ok {
			memberMap[info.Name] = true
		} else {
			// duplicated name error
			log.Println("duplicated name failed: ", info.Name)
			return
		}
	}

	mInfos := []model.MemberInfo{}

	for _, info := range originalInfo {
		if !isTargetMember(info.Name) {
			continue
		}
		birth, _ := time.Parse("2006/01/02", info.Birth)
		mInfo, err := handleMemberData(info.Name, birth)
		if err != nil {
			log.Println("handle original member data failed: ", err.Error(), ", name: ", info.Name)
			return
		}
		mInfos = insert(mInfos, *mInfo)
	}

	for _, info := range phase1Info {
		if !isTargetMember(info.Name) {
			continue
		}
		birth, _ := time.Parse("2006/01/02", info.Birth)
		mInfo, err := handleMemberData(info.Name, birth)
		if err != nil {
			log.Println("handle phase1 member data failed: ", err.Error(), ", name: ", info.Name)
			return
		}
		mInfos = insert(mInfos, *mInfo)
	}

	for _, info := range phase2Info {
		if !isTargetMember(info.Name) {
			continue
		}
		birth, _ := time.Parse("2006/01/02", info.Birth)
		mInfo, err := handleMemberData(info.Name, birth)
		if err != nil {
			log.Println("handle phase2 member data failed: ", err.Error(), ", name: ", info.Name)
			return
		}
		mInfos = insert(mInfos, *mInfo)
	}

	for _, info := range waitingInfo {
		if !isTargetMember(info.Name) {
			continue
		}
		birth, _ := time.Parse("2006/01/02", info.Birth)
		mInfo, err := handleMemberData(info.Name, birth)
		if err != nil {
			log.Println("handle waiting member data failed: ", err.Error(), ", name: ", info.Name)
			return
		}
		mInfos = insert(mInfos, *mInfo)
	}

	results := []model.Result{}
	for i := 0; i < len(mInfos); i++ {
		results = append(results, model.Result{
			SerialNum:        strconv.Itoa(i + 1),
			Name:             mInfos[i].Name,
			Birth:            mInfos[i].Birth.Format("2006/01/02")[1:],
			Age:              formatAge(mInfos[i].Age[0], mInfos[i].Age[1], mInfos[i].Age[2]),
			VerificationType: mInfos[i].VerificationType.String(),
		})
	}

	filename := "anan-file"
	err = genCSVFile(filename, results)
	if err != nil {
		log.Println("gen csv failed: ", err.Error())
		return
	}

	countTypeResult := model.CountTypeResult{}
	for vType, count := range countVerificationType {
		switch vType {
		case constant.VerificationTypeThree:
			countTypeResult.VerificationTypeThree = count
		case constant.VerificationTypeThreeHalf:
			countTypeResult.VerificationTypeThreeHalf = count
		case constant.VerificationTypeFour:
			countTypeResult.VerificationTypeFour = count
		case constant.VerificationTypeFive:
			countTypeResult.VerificationTypeFive = count
		}
	}

	filenameCount := "anan-file-count"
	err = genCSVFile(filenameCount, []model.CountTypeResult{countTypeResult})
	if err != nil {
		log.Println("gen csv failed: ", err.Error())
		return
	}

	log.Printf("gen %s success\n", filename)
}

func insert(mInfos []model.MemberInfo, mInfo model.MemberInfo) []model.MemberInfo {
	mInfos = append(mInfos, mInfo)
	sort(mInfos)
	return mInfos
}

func sort(mInfos []model.MemberInfo) {
	for i := len(mInfos) - 2; i >= 0; i-- {
		if isAgeGreaterThan(mInfos[i], mInfos[i+1]) {
			mInfos[i], mInfos[i+1] = mInfos[i+1], mInfos[i]
		}
	}
}

func isAgeGreaterThan(info1 model.MemberInfo, info2 model.MemberInfo) bool {
	for i := 0; i < len(info1.Age); i++ {
		switch {
		case info1.Age[i] > info2.Age[i]:
			return true
		case info1.Age[i] < info2.Age[i]:
			return false
		}
	}
	return true
}

func genCSVFile(filename string, data interface{}) error {

	file, err := os.OpenFile(filename+".csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println("open file failed: ", err.Error())
		return err
	}
	defer file.Close()

	err = gocsv.MarshalFile(data, file)
	if err != nil {
		log.Println("marshal csv failed: ", err.Error())
		return err
	}

	return nil
}

func isTargetMember(name string) bool {
	return memberMap[name]
}

func handleMemberData(name string, birth time.Time) (*model.MemberInfo, error) {
	ageInts, err := getAge(targetDate, birth)
	if len(ageInts) < 3 {
		log.Println("calculate age failed")
		return nil, err
	}

	vType := getVerificationType(ageInts[0], ageInts[1], ageInts[2])

	// count verification type
	countVerificationType[vType]++

	// to 民國年
	rocBirth := time.Date(birth.Year()-1911, birth.Month(), birth.Day(), 0, 0, 0, 0, location)

	return &model.MemberInfo{
		Name:             name,
		Birth:            rocBirth,
		Age:              ageInts,
		VerificationType: vType,
	}, nil
}

func getVerificationType(y int, m int, d int) constant.VerificationType {
	switch {
	case y == 5:
		return constant.VerificationTypeFive
	case y == 4:
		return constant.VerificationTypeFour
	case y == 3:
		if m > 6 {
			return constant.VerificationTypeThreeHalf
		} else {
			return constant.VerificationTypeThree
		}
	default:
		return constant.VerificationTypeUnkonwn
	}
}

func getAge(cur time.Time, birth time.Time) ([]int, error) {

	// implement curBuf - birthBuf
	curBuf := util.NewCalNumber(
		util.NewNumber(cur.Year(), 10),
		util.NewNumber(int(cur.Month()), 12),
		util.NewNumber(cur.Day(), constant.MonthDays(cur.Month())),
	)
	birthBuf := util.NewCalNumber(
		util.NewNumber(birth.Year(), 10),
		util.NewNumber(int(birth.Month()), 12),
		util.NewNumber(birth.Day(), constant.MonthDays(birth.Month())),
	)

	diff, err := curBuf.Sub(birthBuf)
	if err != nil {
		return nil, err
	}

	return diff.ToInts(), nil
}

func formatAge(y int, m int, d int) string {
	return fmt.Sprintf("%d年%d月%d天", y, m, d)
}

func getPath(fileName constant.FileName, fileType constant.FileType, dirs ...string) string {
	dirPath := strings.Join(dirs, "/")
	return path.Join(dirPath, fmt.Sprintf("%s.%s", fileName, fileType))
}

func getMemberFromCSV(filePath string) ([]model.Member, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("get member from csv failed: ", err.Error())
		return nil, err
	}

	defer file.Close()

	members := []model.Member{}
	if err := gocsv.UnmarshalFile(file, &members); err != nil {
		return nil, err
	}

	return members, nil
}

func getPhase1FromCSV(filePath string) ([]model.Phase1, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("get phase1 from csv failed: ", err.Error())
		return nil, err
	}

	defer file.Close()

	phase1s := []model.Phase1{}
	if err := gocsv.UnmarshalFile(file, &phase1s); err != nil {
		return nil, err
	}

	return phase1s, nil
}

func getPhase2FromCSV(filePath string) ([]model.Phase2, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("get phase2 from csv failed: ", err.Error())
		return nil, err
	}

	defer file.Close()

	phase2s := []model.Phase2{}
	if err := gocsv.UnmarshalFile(file, &phase2s); err != nil {
		return nil, err
	}

	return phase2s, nil
}

func getWaitingFromCSV(filePath string) ([]model.Waiting, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("get waiting from csv failed: ", err.Error())
		return nil, err
	}

	defer file.Close()

	waiting := []model.Waiting{}
	if err := gocsv.UnmarshalFile(file, &waiting); err != nil {
		return nil, err
	}

	return waiting, nil
}

func getOriginalFromCSV(filePath string) ([]model.Original, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("get original from csv failed: ", err.Error())
		return nil, err
	}

	defer file.Close()

	originals := []model.Original{}
	if err := gocsv.UnmarshalFile(file, &originals); err != nil {
		return nil, err
	}

	return originals, nil
}
