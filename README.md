# pmimport
   - 媒体导入工具(Personal Media Import)
   - 作为照片视频管理系统的文件导入工具

# 照片视频管理

- 基本设计理念： 文件的存储要避免供应商锁定（不依赖于特定厂商的专用工具和技术），使得即使用户没有本系统以能正常访问存储的文件，以及可以用本系统定义的存储规范检索到文件，基于这个原则，对于文件的存储使用标准的操作系统的文件系统，存储的安全性可用磁盘阵列、跨设备的目录同步，Ceph等解决方案．

## 存储规范

```shell
# 磁盘目录结构
storge_dir/                #存储根目录
  + username /             #用户目录
     + logs/
       - improt_xxxxxx.log #每次导入的日志文件，以时间为文件名
     + media /             #媒体文件目录
       + yyyy /            #年份目录，一年一个目录
         + yyyy-mm-dd /       #日期目录，一天一个目录
           - photoFileSha1.jpg         # 照片文件，文件SHA1作为文件名, 扩展名使用原文件的扩展名，如.jpg,jpeg,raw,mp4.mov等
           - photoFileSha1.thumb.jpg   # 对应文件的缩略图，以thumb.jpg结尾
           - photoFileSha1_info.json   # 对应文件的信息文件json格式, 保存文件的拍摄信息，标签等
           
     + photo_album         #相册目录
         + 相册1           #相册名
            + 相册1.album.idx    #相册的索引文件， 记录了本相册对应的照片的路径
            + 相册1.secert       #相册的权限信息
```

### 文件说明

1. photoFileSha1_info.json 文件

   - 信息文件使用json格式保存对应媒体文件的信息，如拍摄相机，镜头，用户标签， GPS等

   - 格式如下：

     ```json
     {
         "photoFileSha1.jpg": {     //以文件名作为对象名
             "create_time":"string",
             "Model": "string",         //相机型号
             "LensModel": "string",     //相机镜头
             "LatLong":"string",        //拍摄坐标
             "tags":"tag1, tag2",       //不能包含特殊字符
             "album_text": "string",    //在相册展示时的说明
             "remark":"string"          //文件的笔记，对特殊字符进行转义
             //文件照片文件中提取的其它信息字段
         }
     }
     ```

2. xxx.album.idx 相册索引文件

   - 相册文件描述了本相册的基本信息，如相册故事，每个文件的说明

     ```json
     //相册索引文件
     {
         "thestory":"string",   //相册故事
         "create_time":"string",
         "update_time":"string",
         "files": [
             {
                 "file":"filepath",
                 "album_text": "string",   //覆盖原文件中的"album_text"
             }
         ]
     }
     ```

     

## 基本功能

### 导入程序
- 指定用户，将文件导入指定用户的目录下
- 导入指定目录，或指定的文件到存储系统中
- 导入的时候可以为文件指定一个或多个用户标签(tag)
- 可以指定重复文件的导入规则，忽略，覆盖，重命名， 并打印重复文件列表
- 可以指定源文件的处理规则，删除，改名

### 配置文件
- 配置文件pmi.yaml

  ```yaml
  storage:
    path: /opt/storage
    userid: liujiaxiang
  ```

  注: 

  - storage.path: 文件仓库的根目录路径
  - storage.userid: 当前导入操作的用户的id

##　基本使用

- 导入指定目录下的文件到媒体库，重复的文件不做操作

  ```bash
  pmimport import --tags "test_photo;tags2" --from=~/files/
  ```

- 导入指定目录下的文件到媒体库，有重复的文件时，将仓库中的旧文件改名在导入

  ```bash
  pmimport import --tags "test_photo;tags2"  --overwrite  --from=~/files/
  ```

- 导入指定目录下的文件到媒体库，导入成功后将源文件改为，加上"import_"前缀

  ```bash
  pmimport import --tags "test_photo;tags2"  --rename  --from=~/files/
  ```

- 导入指定目录下的文件到媒体库，导入成功后将源文件删除   <u>***谨慎使用!!!***</u>

  ```bash
  # 谨慎使用 !!!
  pmimport import --tags "test_photo;tags2"  --destroy  --from=~/files/
  ```

- 导入时，排除指定文件名前缀的文件

  ```bash
  pmimport import --tags "test_photo;tags2"  --exclude-file="VID_"  --from=~/files/
  ```

- 当照片中不能正确提取相机信息时，指定相机信息

  ```bash
  pmimport import --tags "test_photo;tags2"  --model="Canon"  --from=~/files/
  ```

  

