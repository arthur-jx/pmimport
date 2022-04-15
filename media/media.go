package media

import "time"

//每个媒体文件对应信息
type MediaFileInfo struct {
	FileHash   string    `json:"file_hash"`   //文件的ｈａｓｈ串
	CreateTime time.Time `json:"create_time"` //照片拍摄时间
	Model      string    `json:"model"`       //相机型号
	LensModel  string    `json:"lens_model"`  //相机镜头
	LatLong    string    `json:"latlong"`     //拍摄坐标
	Tags       string    `json:"tags"`        //标签，用分号分隔，不能包含特殊字符
	AlbumText  string    `json:"album_text"`  //在相册展示时的说明
	Remark     string    `json:"remark"`      //文件的笔记，对特殊字符进行转义
}
