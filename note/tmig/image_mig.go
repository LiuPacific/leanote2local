package tmig

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"leanote2vnote/note/constants"
	"leanote2vnote/note/model"
	"os"
	"regexp"
	"strings"
)

var LeanoteImageLinkRegStr = "!\\[.*?\\]\\(/api/file/getImage\\?fileId=(\\w*?)\\)"
var LeanoteImageLinkRegCompiled = regexp.MustCompile(LeanoteImageLinkRegStr)

//replace ![title](/api/file/getImage?fileId=5f1d3aba4b6c3f6fa6000e49)
//to      ![title](_tp_images/20200905232545067_24522.png)
func ReplaceImageLink(session *mgo.Session, notebookPath string, content string) (contentResult string) {
	submatchs := LeanoteImageLinkRegCompiled.FindAllStringSubmatch(content, -1)
	for i := 0; i < len(submatchs); i++ {
		mdLink := CopyImageLink(session, notebookPath, submatchs[i][1])

		content = strings.Replace(content, submatchs[i][0], mdLink, 1)
	}
	return content
}

//replace ![title](/api/file/getImage?fileId=5f1d3aba4b6c3f6fa6000e49)
//to      ![title](_tp_images/20200905232545067_24522.png)
func CopyImageLink(session *mgo.Session, notebookPath string, imageFileId string) (mdLink string) {
	//![title](/api/file/getImage?fileId=5f1d3aba4b6c3f6fa6000e49)
	//
	//> db.files.find({"_id":ObjectId("5f1d3aba4b6c3f6fa6000e49")}).sort({_id:-1}).pretty().limit(3)
	//{
	//	"_id" : ObjectId("5f1d3aba4b6c3f6fa6000e49"),
	//	"UserId" : ObjectId("5d5cecedf679617402000001"),
	//	"AlbumId" : ObjectId("52d3e8ac99c37b7f0d000001"),
	//	"Name" : "eed28224d3543443a249b9e3b5405f75.png",
	//	"Title" : "UnTitled",
	//	"Size" : NumberLong(10210),
	//	"Type" : "",
	//	"Path" : "files/571/5d5cecedf679617402000001/99/images/eed28224d3543443a249b9e3b5405f75.png",
	//	"IsDefaultAlbum" : true,
	//	"CreatedTime" : ISODate("2020-07-26T08:11:38.258Z")
	//}

	collectionFile := session.DB(model.DbLeanoteName).C(model.CollectionFilesName)
	result := model.File{}
	err := collectionFile.
		Find(bson.M{"_id": bson.ObjectIdHex(imageFileId)}).
		One(&result)
	if err != nil {
		fmt.Printf("try find record error[%s]\n", err.Error())
		panic(err)
	}

	relativeImagePath := constants.DefaultImageFolder + "/" + result.Name
	mdLink = fmt.Sprintf("![title](%s)", relativeImagePath) //![title](_tp_images/20200905232545067_24522.png)

	targetImagePath := notebookPath + "\\" + relativeImagePath
	sourceImagePath := constants.LeanoteFileDir + "\\" + result.Path

	copy(sourceImagePath, targetImagePath)

	return mdLink
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
