package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

//{
//    "created_time": "2020-09-05T15:19:03Z",
//    "files": [
//        {
//            "attachment_folder": "",
//            "attachments": [
//            ],
//            "created_time": "2020-09-05T15:19:17Z",
//            "modified_time": "2020-09-05T15:25:47Z",
//            "name": "name0.md",
//            "tags": [
//            ]
//        }
//    ],
//    "sub_directories": [
//        {
//            "name": "liberation_son"
//        }
//    ],
//    "version": "1"
//}
func TestVNoteJson(t *testing.T) {
	vInfo := &VNotebookInfo{
		CreatedTime: "2020-09-06T00:00:00Z",
		Files: []VNoteInfo{
			VNoteInfo{
				AttachmentFolder: "",
				Attachments:      make([]string, 0, 0),
				CreatedTime:      "2020-09-06T00:00:00Z",
				ModifiedTime:     "2020-09-06T00:00:00Z",
				Name:             "name0.md",
				Tags:             []string{},
			},
		},
		SubDirectories: []SubDirInfo{
			SubDirInfo{Name: "liberation_son"},
		},
		Version: "1",
	}

	marshal, err := json.Marshal(vInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(marshal))
}
