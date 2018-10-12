package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
)

func main() {
	//設定ファイル読み込み
	confFile, err := ioutil.ReadFile("C:\\dev\\credentials\\Vision.json")
	if err != nil {
		log.Fatalln("failed to read configuration file", err)
	}

	// コマンドライン引数で画像URL取得
	flag.Parse()
	strOpt := flag.Arg(0)

	//画像ファイル読み込み
	localPath, err := ImgFileDownload(strOpt)
	imgData, err := ioutil.ReadFile(localPath)
	if err != nil {
		log.Fatalln("failed to read image file", err)
	}

	cfg, err := google.JWTConfigFromJSON([]byte(confFile), vision.CloudPlatformScope)
	client := cfg.Client(context.Background())

	svc, err := vision.New(client)
	enc := base64.StdEncoding.EncodeToString([]byte(imgData))
	img := &vision.Image{Content: enc}

	feature := &vision.Feature{
		//Type:       "WEB_DETECTION",
		Type:       "SAFE_SEARCH_DETECTION",
		MaxResults: 3,
	}

	req := &vision.AnnotateImageRequest{
		Image:    img,
		Features: []*vision.Feature{feature},
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}
	res, err := svc.Images.Annotate(batch).Do()

	body, err := json.MarshalIndent(res.Responses[0], "", "\t")
	fmt.Println(string(body))
}

// ImgFileDownload はURLから画像ダウンロードする
func ImgFileDownload(url string) (path string, err error) {
	imgFilePath := "imgs/img_" + getNow() + ".jpg"

	response, err := http.Get(url)
	if err != nil {
		log.Fatalln("failed to get image file", err)
	}

	defer response.Body.Close()

	file, err := os.Create(imgFilePath)
	if err != nil {
		log.Fatalln("failed to create image file", err)
	}
	defer file.Close()

	io.Copy(file, response.Body)
	return imgFilePath, err
}

// getNow は現在時刻を文字列で返す
func getNow() string {
	n := time.Now()
	return n.Format("150405.000")
}
