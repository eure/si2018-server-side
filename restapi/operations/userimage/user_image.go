package userimage

import (
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/eure/si2018-server-side/repositories"

	"github.com/eure/si2018-server-side/entities"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func xorshift(n uint64) uint64 {
	n = n ^ (n << 7)
	return n ^ (n >> 9)
}

func quickEncrypt(n uint64) uint64 {
	n ^= 8123456789123456789
	for i := 0; i < 10; i++ {
		n = xorshift(n)
	}
	return n
}

func getFreeFileName(path, extension string) (string, error) {
	// 時刻がバレないように簡易的に暗号化したものをファイル名にする
	code := quickEncrypt(uint64(time.Now().UnixNano()))
	// 同じ時刻では 100 枚までしか写真を保存しない
	// 1 ns 単位で重複することはほぼないはずなので, サーバーの時刻を止めない限り問題ない
	const limit = 100
	format := "img%x%02d" + extension
	for i := 0; i < limit; i++ {
		filename := fmt.Sprintf(format, code, i)
		if !exists(path + filename) {
			return filename, nil
		}
	}
	return "", errors.New("storage capacity exceeded")
}

// image のタグが tag と一致するか
func tagMatches(tag []byte, image strfmt.Base64) bool {
	seq := []byte(image)
	if len(seq) < len(tag) {
		return false
	}
	for i, t := range tag {
		if seq[i] != t {
			return false
		}
	}
	return true
}

var pngTag = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

// image のファイル形式が png であるか
func isPngFile(image strfmt.Base64) bool {
	return tagMatches(pngTag, image)
}

var jpgTag = []byte{0xff, 0xd8}

// image のファイル形式が jpg であるか
func isJpgFile(image strfmt.Base64) bool {
	return tagMatches(jpgTag, image)
}

var gifTag = []byte{0x47, 0x49, 0x46, 0x38}

// image のファイル形式が gif であるか
func isGifFile(image strfmt.Base64) bool {
	return tagMatches(gifTag, image)
}

// image が png, jpg, gif 形式であれば対応する拡張子 ".png" などを返す
func getFileType(image strfmt.Base64) (string, error) {
	if isPngFile(image) {
		return ".png", nil
	}
	if isJpgFile(image) {
		return ".jpg", nil
	}
	if isGifFile(image) {
		return ".gif", nil
	}
	return "", errors.New("unknown file type")
}

// 画像サイズの制限
// ギリギリ 10.1 MB 未満にすることで, 10 MB に制限
const ImageSizeLimit = 10590617

// DB アクセス: 3 回
// 計算量: O(1)
func PostImage(p si.PostImagesParams) middleware.Responder {
	if p.Params.Token == "" {
		return si.PostImageThrowBadRequest("missing token")
	}
	if p.Params.Image == nil {
		return si.PostImageThrowBadRequest("broken image")
	}
	var err error
	// 画像のバリデーション
	if len(p.Params.Image) > ImageSizeLimit {
		return si.PostImageThrowBadRequest("Image is too large")
	}
	extension, err := getFileType(p.Params.Image)
	if err != nil {
		return si.PostImageThrowInternalServerError("getFileType", err)
	}
	imageRepo := repositories.NewUserImageRepository()
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Params.Token)
		if err != nil {
			return si.PostImageThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return si.PostImageThrowUnauthorized("GetByToken failed")
		}
		id = token.UserID
	}
	// 既にある画像のタイムスタンプを使いまわしたい
	image, err := imageRepo.GetByUserID(id)
	if err != nil {
		return si.PostImageThrowInternalServerError("GetByUserID", err)
	}
	var localFile, serverFile string
	{
		path := os.Getenv("ASSETS_PATH")
		uri := os.Getenv("ASSETS_BASE_URI")
		filename, err := getFreeFileName(path, extension)
		if err != nil {
			return si.PostImageThrowInternalServerError("getFreeFileName", err)
		}
		localFile = path + filename
		serverFile = uri + filename
	}
	// 画像をローカルに保存する
	{
		file, err := os.Create(localFile)
		if err != nil {
			return si.PostImageThrowInternalServerError("Create", err)
		}
		fmt.Println("file:", localFile, ", size:", len(p.Params.Image))
		// file.Close() でも error が吐かれる可能性がある
		// とりあえず今は対処しない
		defer file.Close()
		file.Write(p.Params.Image)
	}
	// DB を更新
	now := strfmt.DateTime(time.Now())
	ent := entities.UserImage{
		UserID:    id,
		Path:      serverFile,
		UpdatedAt: now}
	// 既に DB にデータがあれば Update、なければ Create
	if image == nil {
		err = imageRepo.Create(ent)
		if err != nil {
			return si.PostImageThrowInternalServerError("Create", err)
		}
	} else {
		err = imageRepo.Update(ent)
		if err != nil {
			return si.PostImageThrowInternalServerError("Update", err)
		}
	}
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(ent.Path)})
}
