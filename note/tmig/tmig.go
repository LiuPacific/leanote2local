package tmig

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"leanote2vnote/note/common"
	"leanote2vnote/note/constants"
	"leanote2vnote/note/model"
	"os"
	"strconv"
	"strings"
)

func StartMig(session *mgo.Session) error {

	// foreach collection of 'notes'
	// gen folder of the note
	// gen the file of the note
	// fill the note with the content from collection of 'note_contents'
	// handle images
	// -- save it
	// -- modify image link to link it

	collectionNote := session.DB(model.DbLeanoteName).C(model.CollectionNotesName)

	results := []model.Note{}

	err := collectionNote.Find(bson.M{}).All(&results)
	if err != nil {
		fmt.Printf("try find record error[%s]\n", err.Error())
		return err
	}
	fmt.Println(len(results))

	return nil
}

// CreateDir
// - notebook
// -- _tp_images
// -- _vnote.json
// -- sub_notebook
func CreateTopDir(session *mgo.Session) error {
	fmt.Println("start")

	// from top : notebooks  "ParentNotebookId" : "" and IsDeleted = true
	collectionNotebooks := session.DB(model.DbLeanoteName).C(model.CollectionNotebooksName)

	//db.notebooks.find(
	//	{$and:
	//	[
	//		{$or:[{"ParentNotebookId":""},{"ParentNotebookId" :{$exists:false} }]},
	//		{"IsDeleted":false}
	//	]}
	//).sort({_id:-1}).pretty().limit(500)

	topNotebooksResults := make([]*model.Notebook, 0)
	err := collectionNotebooks.
		Find(bson.M{
			"$and": []bson.M{
				{"IsDeleted": false},
				{"$or": []bson.M{
					{"ParentNotebookId": nil}, //unset {"ParentNotebookId": "{$exists:false}"},
					{"ParentNotebookId": ""},
				}},
				{"UserId": bson.ObjectIdHex(constants.UserId)},
			},
		}).
		All(&topNotebooksResults)

	fmt.Println("len is " + strconv.Itoa(len(topNotebooksResults)))
	if err != nil {
		fmt.Printf("CreateTopDir error[%s]\n", err.Error())
		return err
	}
	if len(topNotebooksResults) == 0 {
		return errors.New("top notebook length = 0")
	}

	topVInfo := &model.VNotebookInfo{
		AttachmentFolder: constants.DefaultAttachmentFolder,
		CreatedTime:      constants.DefaultTime,
		Files:            []model.VNoteInfo{},
		ImageFolder:      constants.DefaultImageFolder,
		RecycleBinFolder: "_v_recycle_bin",
		SubDirectories: []model.SubDirInfo{
			model.SubDirInfo{Name: "liberation"},
		},
		Tags:    []string{},
		Version: "1",
	}

	for _, topNotebook := range topNotebooksResults {
		fmt.Println("- topNotebook " + topNotebook.Title)
		notebookTitle := topNotebook.Title
		notebookDirPath := constants.VNoteRootDir + "//" + notebookTitle
		err := CreateNotebooks(session, notebookDirPath, topNotebook.NotebookId)
		if err != nil {
			return err
		}
		topVInfo.SubDirectories = append(topVInfo.SubDirectories, model.SubDirInfo{Name: topNotebook.Title})
	}

	topVNoteJson, err := json.Marshal(topVInfo)
	if err != nil {
		return err
	}

	topVNodeJsonPath := constants.VNoteRootDir + "//_vnote.json"
	err = common.ReWriteToFile(topVNodeJsonPath, topVNoteJson)
	if err != nil {
		return err
	}

	fmt.Println("end")
	return nil
}

// include note and notebook
// include generating notes and _vnote.json
func CreateNotebooks(session *mgo.Session, notebookPath string, notebookId bson.ObjectId) error {
	vnoteInfos := make([]model.VNoteInfo, 0)
	subDirInfos := make([]model.SubDirInfo, 0)

	// generate notebook dir
	err := os.MkdirAll(notebookPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(notebookPath+"\\"+constants.DefaultImageFolder, os.ModePerm)
	if err != nil {
		return err
	}

	// generate markdown files
	noteResults := make([]*model.Note, 0)
	collectionNotes := session.DB(model.DbLeanoteName).C(model.CollectionNotesName)
	err = collectionNotes.Find(bson.M{
		"$and": []bson.M{
			{"NotebookId": notebookId},
			{"IsDeleted": false},
			{"IsMarkdown": true},
		},
	}).All(&noteResults)
	if err != nil {
		return err
	}
	for _, note := range noteResults {
		fmt.Println("-- note " + note.Title)
		windowsIllegalTokens := []string{"?", ",", "\\", "/", "*", "\"", "“", "”", "<", ">", "|"}
		//note.Title = strings.ReplaceAll(note.Title, "\\", "-")
		//note.Title = strings.ReplaceAll(note.Title, "/", "-")
		for _, str := range windowsIllegalTokens {
			note.Title = strings.ReplaceAll(note.Title, str, "-")
		}
		note.Title = note.Title + ".md"
		vnoteInfos = append(vnoteInfos, model.VNoteInfo{
			AttachmentFolder: constants.DefaultAttachmentFolder,
			Attachments:      []string{},
			CreatedTime:      constants.DefaultTime,
			ModifiedTime:     constants.DefaultTime,
			Name:             note.Title,
			Tags:             []string{},
		})
		err = CreateNoteMarkdown(session, note, notebookPath)
		if err != nil {
			return err
		}
	}

	// call CreateNotebooks
	notebookResults := make([]*model.Notebook, 0)
	collectionNotebooks := session.DB(model.DbLeanoteName).C(model.CollectionNotebooksName)
	err = collectionNotebooks.Find(bson.M{
		"$and": []bson.M{
			{"IsDeleted": false},
			{"ParentNotebookId": notebookId},
			{"UserId": bson.ObjectIdHex(constants.UserId)},
		},
	}).All(&notebookResults)
	if err != nil {
		return err
	}
	for _, subNotebook := range notebookResults {
		fmt.Println("-- subNotebook " + subNotebook.Title)
		subDirInfos = append(subDirInfos, model.SubDirInfo{Name: subNotebook.Title})
		CreateNotebooks(session, notebookPath+"//"+subNotebook.Title, subNotebook.NotebookId)
	}

	// create _vnote.json
	vInfo := &model.VNotebookInfo{
		//AttachmentFolder: "tp_attachments",
		CreatedTime: "2020-09-06T00:00:00Z",
		Files:       vnoteInfos,
		//DefaultImageFolder:      constants.DefaultImageFolder,
		//RecycleBinFolder: "_v_recycle_bin",
		SubDirectories: subDirInfos,
		Tags:           []string{},
		Version:        "1",
	}
	vInfoContent, err := json.Marshal(vInfo)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	common.ReWriteToFile(notebookPath+"//_vnote.json", vInfoContent)
	return nil
}

func CreateNoteMarkdown(session *mgo.Session, note *model.Note, notebookPath string) error {
	collectionNoteContents := session.DB(model.DbLeanoteName).C(model.CollectionNoteContentsName)

	noteContentResult := &model.NoteContent{}
	err := collectionNoteContents.Find(bson.M{
		"$and": []bson.M{
			{"UserId": bson.ObjectIdHex(constants.UserId)},
			{"_id": note.NoteId},
			//{"IsDeleted": false},
			//{"IsTrash": false},
		},
	}).One(noteContentResult)
	if err != nil {
		fmt.Println(note.NoteId.String())
		return err
	}

	//replace ![title](/api/file/getImage?fileId=5f1d3aba4b6c3f6fa6000e49)
	//to      ![title](_tp_images/20200905232545067_24522.png)
	noteContentResult.Content = ReplaceImageLink(session, notebookPath, noteContentResult.Content)

	common.ReWriteToFile(notebookPath+"//"+note.Title, ([]byte)(noteContentResult.Content))

	return nil
}
