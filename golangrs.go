package golangrs

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// >For a Timer created with NewTimer, Reset should be invoked only on stopped or expired timers with drained channels.
func SafelyResetTimer(t *time.Timer, d time.Duration) { //note this is NOT really safe if you called Stop() already or read the channel already
	if !t.Stop() {
		<-t.C
	}
	t.Reset(d)
}
func StopOrDrainTimer(t *time.Timer) { //note this should help garbage collection because otherwise GC does not touch it until it expires
	if !t.Stop() {
		<-t.C
	}
}

func DoDeleteRequest(hc *http.Client, url string) (int, []byte) {
	req, err := http.NewRequest("DELETE", url, nil)
	CheckErr(err)
	resp, err := hc.Do(req)
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body) //note in order to reuse http connection, you do 2 thing: read all bytes and close body!
	resp.Body.Close()                      //here omits checking close error
	CheckErr(err)
	return resp.StatusCode, body
}
func DoGetRequestWithUserAgentFor2xx(hc *http.Client, url, userAgent string) []byte {
	cod, body := DoGetRequestWithUserAgent(hc, url, userAgent)
	if cod > 299 || cod < 200 {
		CheckErrLogMsg(`HTTP` + strconv.Itoa(cod))
	}
	return body
}
func DoGetRequestWithHeaders(hc *http.Client, url string, headersMap map[string]string) (int, []byte) {
	req, err := http.NewRequest("GET", url, nil)
	CheckErr(err)
	for k, v := range headersMap {
		req.Header.Set(k, v)
	}
	resp, err := hc.Do(req)
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body) //note in order to reuse http connection, you do 2 thing: read all bytes and close body!
	resp.Body.Close()                      //here omits checking close error
	CheckErr(err)
	return resp.StatusCode, body
}
func DoGetRequestWithUserAgent(hc *http.Client, url, userAgent string) (int, []byte) {
	req, err := http.NewRequest("GET", url, nil)
	CheckErr(err)
	req.Header.Set("User-Agent", userAgent)
	resp, err := hc.Do(req)
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body) //note in order to reuse http connection, you do 2 thing: read all bytes and close body!
	resp.Body.Close()                      //here omits checking close error
	CheckErr(err)
	return resp.StatusCode, body
}
func DoGetRequestFor200(hc *http.Client, url string) []byte {
	cod, body := DoGetRequest(hc, url)
	if cod != 200 {
		CheckErrLogMsg(`HTTP` + strconv.Itoa(cod))
	}
	return body
}
func DoGetRequest(hc *http.Client, url string) (int, []byte) {
	resp, err := hc.Get(url)
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body) //note in order to reuse http connection, you do 2 thing: read all bytes and close body!
	resp.Body.Close()                      //here omits checking close error
	CheckErr(err)
	return resp.StatusCode, body
}
func DoPostRequestFor200(hc *http.Client, url, contentType, strAsBody string) []byte {
	cod, body := DoPostRequest(hc, url, contentType, strAsBody)
	if cod != 200 {
		CheckErrLogMsg(`HTTP` + strconv.Itoa(cod))
	}
	return body
}
func DoPostRequest(hc *http.Client, url, contentType, strAsBody string) (int, []byte) {
	resp, err := hc.Post(url, contentType, strings.NewReader(strAsBody))
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body) //note in order to reuse http connection, you do 2 thing: read all bytes and close body!
	resp.Body.Close()                      //here omits checking close error
	CheckErr(err)
	return resp.StatusCode, body
}
func DoPostRequestWithoutBody(hc *http.Client, url, contentType string) (int, []byte) {
	resp, err := hc.Post(url, contentType, nil)
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body) //note in order to reuse http connection, you do 2 thing: read all bytes and close body!
	resp.Body.Close()                      //here omits checking close error
	CheckErr(err)
	return resp.StatusCode, body
}

func LogIntSliceValuesNIndexFromMaxToMin(ints []int, discardThreshold int) {
	dup := make([]int, len(ints))
	copy(dup, ints)
	sort.Ints(dup)
	var lastIntSet bool
	var lastInt int
	for i := len(dup) - 1; i >= 0; i-- {
		nextInt := dup[i]
		if nextInt <= discardThreshold {
			return
		}
		if lastIntSet {
			if nextInt == lastInt {
				continue
			}
		} else {
			lastIntSet = true
		}
		lastInt = nextInt
		for ind, v := range ints {
			if v == nextInt {
				log.Println(ind, v)
			}
		}
	}
}

func CreateSqliteSimpleLogTbl(db *sql.DB) {
	DbExecSql(db, `create table IF NOT EXISTS simpleLog(s1 text not null, ms1970 integer not null)`)
}
func CheckSqliteTblExists(db *sql.DB, tblName string) bool {
	cnt := DbQueryRowGetInt(db, `select count(*) from sqlite_schema where lower(tbl_name)=?`, strings.ToLower(tblName))
	return cnt != 0
}
func TxExecSql(tx *sql.Tx, query string, args ...interface{}) sql.Result {
	retval, err := tx.Exec(query, args...)
	CheckErr(err)
	return retval
}
func DbExecSql(db *sql.DB, query string, args ...interface{}) sql.Result {
	retval, err := db.Exec(query, args...)
	CheckErr(err)
	return retval
}
func DbQueryRowGetStr(db *sql.DB, query string, args ...interface{}) string {
	var retval string
	err := db.QueryRow(query, args...).Scan(&retval)
	CheckErr(err)
	return retval
}
func DbQueryRowGetInt64(db *sql.DB, query string, args ...interface{}) int64 {
	var retval int64
	err := db.QueryRow(query, args...).Scan(&retval)
	CheckErr(err)
	return retval
}
func DbQueryRowGetInt(db *sql.DB, query string, args ...interface{}) int { //? is int really ok? mattn sqlite docs has example of using int, but database/sql docs always use int64?
	var retval int
	err := db.QueryRow(query, args...).Scan(&retval)
	CheckErr(err)
	return retval
}
func DbQueryRowAtLeastOneExists(db *sql.DB, query string, args ...interface{}) bool {
	rows, err := db.Query(query, args...)
	CheckErr(err)
	defer rows.Close()
	retval := rows.Next()
	CheckErr(rows.Err())
	return retval
}
func DbQueryIntoInt64Slice(db *sql.DB, query string, args ...interface{}) (retval []int64) {
	rows, err := db.Query(query, args...)
	CheckErr(err)
	defer rows.Close()
	var outValue int64
	for rows.Next() {
		err := rows.Scan(&outValue)
		CheckErr(err)
		retval = append(retval, outValue)
	}
	CheckErr(rows.Err())
	return
}

func WriteShortStrToBytesBuffer(buf *bytes.Buffer, str string) {
	//WriteBlobToBytesBuffer(buf, []byte(str))//note WriteBlobToBytesBuffer is slower!! because []byte(str) allocates an extra mem object!!
	blen := len(str)
	if blen > 127 {
		CheckErrLogMsg(`Long str cannot be written as short str`)
	}
	buf.WriteByte(uint8(blen))
	buf.WriteString(str)
}
func WriteBlobToBytesBuffer(buf *bytes.Buffer, blob []byte) {
	blen := len(blob)
	if blen < 255 {
		buf.WriteByte(uint8(blen))
	} else {
		buf.WriteByte(255)
		WriteUint32ToBytesBuffer(buf, uint32(blen))
	}
	buf.Write(blob)
}
func WriteStrToBytesBuffer(buf *bytes.Buffer, str string) {
	//WriteBlobToBytesBuffer(buf, []byte(str))//note WriteBlobToBytesBuffer is slower!! because []byte(str) allocates an extra mem object!!
	blen := len(str)
	if blen < 255 {
		buf.WriteByte(uint8(blen))
	} else {
		buf.WriteByte(255)
		WriteUint32ToBytesBuffer(buf, uint32(blen))
	}
	buf.WriteString(str)
}
func WriteUnixMilliAsUint64(buf *bytes.Buffer, gotime time.Time) {
	timeUint := uint64(gotime.UnixMilli())
	WriteUint64ToBytesBuffer(buf, timeUint)
}
func WriteUint32ToBytesBuffer(buf *bytes.Buffer, i uint32) {
	var bytesbuf [4]byte
	binary.LittleEndian.PutUint32(bytesbuf[:], i)
	buf.Write(bytesbuf[:]) //no need to check retval bc "err is always nil"
}
func WriteUint64ToBytesBuffer(buf *bytes.Buffer, i uint64) {
	var bytesbuf [8]byte
	binary.LittleEndian.PutUint64(bytesbuf[:], i)
	buf.Write(bytesbuf[:]) //no need to check retval bc "err is always nil"
}
func WriteInt64ToBytesBuffer(buf *bytes.Buffer, i int64) {
	WriteUint64ToBytesBuffer(buf, uint64(i))
}

func ReadStdinLinesToChan(chOut chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		chOut <- scanner.Text()
	}
	//close(chOut)//closing it really does not server any purpose?
}
func ReadFileLines(filenm string) (retval []string) {
	file, err := os.Open(filenm)
	CheckErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { //note this technique does not load whole file. it loads line by line
		line := scanner.Text()
		retval = append(retval, line)
	}
	err = scanner.Err() //note when scanner.Scan returns false, Err() could indicate an error (when it is just EOF, Err() gives nil)
	CheckErr(err)
	return
}
func ReadHugeFilelines(filenm string, cb func(string) bool) {
	file, err := os.Open(filenm)
	CheckErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { //note this technique does not load whole file. it loads line by line
		line := scanner.Text()
		if !cb(line) {
			break
		}
	}
	err = scanner.Err() //note when scanner.Scan returns false, Err() could indicate an error (when it is just EOF, Err() gives nil)
	CheckErr(err)
}
func ReadHugeFileFirstLine(filenm string) string {
	var retval string
	file, err := os.Open(filenm)
	CheckErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { //note this technique does not load whole file. it loads line by line
		retval = scanner.Text()
		break
	}
	//if err := scanner.Err(); err != nil {//note when scanner.Scan returns false, Err() could indicate an error (when it is just EOF, Err() gives nil)
	//	log.Fatal(err)
	//}
	return retval
}
func FileSizeViaLstat(filenm string) int64 {
	fileInfo, err := os.Lstat(filenm)
	CheckErr(err)
	return fileInfo.Size()
}
func OpenFileForAppendAndRotateIfNecessary(filenm string) *os.File {
	const sizeLimit = 1000000
	fileInfo, err := os.Stat(filenm)
	if err != nil {
		if os.IsNotExist(err) { //>IsNotExist returns a boolean indicating whether the error is known to report that a file or directory does not exist.
		} else {
			CheckErrWhichIsNotNil(err) //panics
		}
	} else if fileInfo.Size() > sizeLimit {
		err = os.Rename(filenm, filenm+`.`+time.Now().Format("20060102150405"))
		CheckErr(err)
	}
	logf, err := os.OpenFile(filenm, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //without O_RDWR, log doesn't work on linux?
	CheckErr(err)
	return logf
}

// ???closing the file or not is quite a problem you cannot decide???
func AssociateFilesWithStdoutStderr(outfn, errfn string) { //note this func is probably useless, bc normally golang always use the standard logger, so the only thing you need is to call log.SetOutput
	if err := os.Stdout.Close(); err != nil {
		panic(err)
	}
	if err := os.Stderr.Close(); err != nil {
		panic(err)
	}
	outf, err := os.OpenFile(outfn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logf, err := os.OpenFile(errfn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer outf.Close() //fixme?? closing the file will cause log never written??
	defer logf.Close()
	os.Stdout = outf
	os.Stderr = logf
}

func CheckPathExists(filenm string) bool {
	_, err := os.Stat(filenm)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		CheckErrWhichIsNotNil(err)
	}
	return true
}
func GetGreatestFilenameOfFileInDir(dirname string) string {
	files, err := os.ReadDir(dirname)
	CheckErr(err)
	for ind := len(files) - 1; ind >= 0; ind-- {
		retval := files[ind]
		if !retval.IsDir() {
			return retval.Name()
		}
	}
	return ""
}
func Float64ToStrUtility(flo float64) string {
	return strconv.FormatFloat(flo, 'f', -1, 64)
}
func Float64ToStrWithPrecision(flo float64, prec int) string {
	return strconv.FormatFloat(flo, 'f', prec, 64)
}
func JsonUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	CheckErr(err)
}

func AssertElsePanicWithStr(ev bool, msg string) {
	if !ev {
		panic(msg)
	}
}
func AssertElsePanic(ev bool) {
	if !ev {
		panic(`Assert failure`)
	}
}
func LogStrAndPanicIfFalse(ok bool, msg string) {
	if !ok {
		log.Println(msg)
		panic(msg)
	}
}

// ? this func is actually quite stupid because you can simple use `debug.Stack()` instead?//This func should print less than debug.Stack() though.//Another benefit is that it looks different from the style of debug.Stack() so when you read log files you easily distinguish immediate err log made by this func and the log made by recover()+debug.Stack().
func logBuilderCallerFn() *strings.Builder {
	//pc, filename, line, _ := runtime.Caller(1)
	//funcName := runtime.FuncForPC(pc).Name()
	//_ = pc
	//log.Println(filename+`:`+strconv.Itoa(line)+`: ERR`, err)
	pc := make([]uintptr, 15) //?should be larger?
	numOfEntries := runtime.Callers(3, pc)
	pc = pc[:numOfEntries]
	frames := runtime.CallersFrames(pc)
	var sb strings.Builder
	sb.WriteString("LOG\n")
	for {
		frame, more := frames.Next()
		sb.WriteByte('\t')
		sb.WriteString(frame.File)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(frame.Line))
		sb.WriteByte(' ')
		sb.WriteString(frame.Function)
		sb.WriteByte('\n')
		if !more {
			break
		}
	}
	return &sb
}
func CheckErrLogMsg(msg string) {
	log.Println(logBuilderCallerFn().String(), msg)
	panic(msg)
}
func CheckErrLikeAssert(ev bool, msg string) {
	if !ev {
		log.Println(logBuilderCallerFn().String(), msg)
		panic(msg)
	}
}
func CheckErrWhichIsNotNil(err error) {
	//note you must add this seemingly useless one more layer of caller because you hard-coded `3` to pass to runtime.Callers()!!
	log.Println(logBuilderCallerFn().String(), err)
	panic(err)
}
func CheckErr(err error) {
	if err != nil {
		log.Println(logBuilderCallerFn().String(), err)
		panic(err)
	}
}
func CheckErrWhichIsNotNilWithoutPanic(err error) {
	log.Println(logBuilderCallerFn().String(), err)
}
func CheckErrWithoutPanic(err error) {
	if err != nil {
		log.Println(logBuilderCallerFn().String(), err)
	}
}

func LogWithStackIfValueIsNotNil(err interface{}) { //note you do not call recover in this func because recover does not work if you have 2 layers of callee. Passing interface is more convenient
	if err != nil {
		log.Println(err, "DEBUG STACK:\n"+string(debug.Stack()))
	}
}

// note MsInDay is a common value for you to use. Actually you have Ms1970, MsInDay, MsInHour, MsInMinute, MsInSecond
func PrintMs1970AsDayHourMinuteSecondOffset(ms1970 int64) {
	fmt.Println(`DAY`, GetIndOfDayFromMs1970(ms1970))
	msInDay := GetMsInDayFromMs1970(ms1970)
	fmt.Println(`HOUR`, GetIndOfHourFromMsInDay(msInDay))
	msInHour := GetMsInHourFromMs1970OrMsInDay(msInDay)
	fmt.Println(`MIN`, GetIndOfMinuteFromMsInHour(msInHour))
	msInMin := GetMsInMinuteFromMs1970OrMsInDayHour(msInHour)
	fmt.Println(msInMin/1000, msInMin%1000)
}

func TruncMs1970ToStartOfDay(ms1970 int64) int64 {
	return ms1970 - ms1970%86400000
}

func GetIndOfDayFromMs1970(ms1970 int64) int64 {
	return ms1970 / 86400000
}
func GetIndOfHourFromMsInDay(msInt int64) int64 {
	return msInt / 3600_000
}
func GetIndOfMinuteFromMsInHour(msInt int64) int64 {
	return msInt / 60_000
}
func GetMsInDayFromMs1970(ms1970 int64) int64 {
	return ms1970 % 86400000
}
func GetMsInHourFromMs1970OrMsInDay(msInt int64) int64 {
	return msInt % 3600_000
}
func GetMsInMinuteFromMs1970OrMsInDayHour(msInt int64) int64 {
	return msInt % 60_000
}

const NumOfMsInDay = 1000 * 3600 * 24 //86400000

type TupleIntIntFloat64 struct {
	I1 int
	I2 int
	F  float64
}
type TupleIntFloat64 struct {
	I1 int
	F1 float64
}
type TupleIntInt struct {
	I1 int
	I2 int
}
type TupleByteIntByteSlice struct {
	B1  byte
	I1  int
	BS1 []byte
}
type TupleStrByteSlice struct {
	S1  string
	BS1 []byte
}

type CircularBufFloat64 struct {
	Cap int
	Buf []float64
	Len int
	Ind int
}

func (cb *CircularBufFloat64) Init(leng int) {
	cb.Cap = leng
	cb.Buf = make([]float64, leng)
}
func (cb *CircularBufFloat64) Add(v float64) {
	cb.Buf[cb.Ind] = v
	cb.Ind++
	if cb.Ind == cb.Cap {
		cb.Ind = 0
	}
	if cb.Len != cb.Cap {
		cb.Len++
	}
}

type Ctx struct {
}

var WgForAllGoroutines sync.WaitGroup
