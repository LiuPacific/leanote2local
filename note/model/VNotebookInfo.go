package model

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
type VNotebookInfo struct {
	AttachmentFolder string       `json:"attachment_folder,omitempty"`
	CreatedTime      string       `json:"created_time,omitempty"`
	Files            []VNoteInfo  `json:"files,omitempty"`
	ImageFolder      string       `json:"image_folder,omitempty"`
	RecycleBinFolder string       `json:"recycle_bin_folder,omitempty"`
	SubDirectories   []SubDirInfo `json:"sub_directories,omitempty"`
	Tags             []string     `json:"tags,omitempty"`
	Version          string       `json:"version,omitempty"`
}

//{
//    "attachment_folder": "",
//    "attachments": [
//    ],
//    "created_time": "2020-09-05T15:19:17Z",
//    "modified_time": "2020-09-05T15:25:47Z",
//    "name": "name0.md",
//    "tags": [
//    ]
//}
type VNoteInfo struct {
	AttachmentFolder string   `json:"attachment_folder,omitempty"`
	Attachments      []string `json:"attachments"`
	CreatedTime      string   `json:"created_time,omitempty"`
	ModifiedTime     string   `json:"modified_time,omitempty"`
	Name             string   `json:"name,omitempty"`
	Tags             []string `json:"tags"`
}

type SubDirInfo struct {
	Name string `json:"name"`
}

//func ()
