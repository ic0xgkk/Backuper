# Backuper 备份助手


## 简介

本工具用于对服务器数据进行自动备份，适用于小型的个人服务器，好比博客等。

## 要求

你需要拥有以下东西：

* GnuPG私钥和公钥，可以自己签发，3分钟即可搞定，此处不再说方法，可以去Google。公钥用于服务端对备份加密，私钥用于本地解密
* 阿里云对象存储，要开好RAM子账户并给好OSS权限，同时生成子账户的AccessKey并新建对象存储桶

## 特别注意

* 所有的上传会直接上传到存储桶的根目录
* 备份的文件名格式为`Backuper-<日期>-<主机名>-<UNIX时间戳>.<扩展名>`
* 每次备份任务执行后共上传两个文件，一个是备份文件的日志（注意不是工具的全局日志），一个是备份的文件（tar.gz经过GnuPG加密后得到的tar.gz.gpg）
* 上传的日志中包含备份文件的消息摘要（SHA256）、错误的文件数和详细的文件日志，其中消息摘要和错误的文件数在日志末尾，建议先检查这里有无问题
* 上传的日志中只包含备份任务相关的日志，如果备份没有成功请检查工作目录下的全局日志`Backuper-GlobalLog-<主机名>.log`
* **所有路径一定要使用绝对路径，一定不要使用相对路径，不然如果出现备份错误或者误删文件不要说我没提醒你**
* **工作目录`work_dir`一定不要设置在内存文件系统中，好比`/tmp`，这可能会导致严重问题**

## 友情提示

备份过多时建议定期清理一下，不然阿里云会对存储进行收费，建议每30天清理一次

## 后续计划

暂且并不打算支持AWS S3等对象存储类型，没有太大的需求。毕竟阿里云也能用，而且有GnuPG加持也不是非常需要担心安全

## 配置方法

保存如下的配置到任意名称的json文件即可

```json
{
    "work_dir": "/data/tmp",
    "pub_key_path": "/data/pubkey.asc",
    "end_point":"",
    "access_key_id": "",
    "access_key_secret":"",
    "bucket_name": "",
    "period_day": 24,
    "start_time_hour": 15,
    "start_time_minute": 30,
    "auto_delete": false,
    "immediate_exec": true,
    "backup_path" :[
        "/etc",
        "/home/www/htdocs",
        "/home/test/test.avi"
    ]
}
```

编码格式建议使用UTF-8，尚未测试中文路径，如果是Linux的话，即便是使用中文路径，确保配置文件编码是和系统编码方式一致的就好，Windows的可能要使用GBK编码（不确定）

其中：

* work_dir 工作目录，**一定要是绝对路径！该目录一定要是存在的，否则只有在任务执行时才会报错导致备份失败。请不要把工作目录设置到内存文件系统中（好比/tmp）**
* pub_key_path  GnuPG公钥（.asc文件）的**绝对路径**
* end_point  对象存储的EndPoint
* access_key_id  阿里云账户的AccessKeyID
* access_key_secret  阿里云账户的AccessKeySecret
* bucket_name  存储桶名称
* period_day  备份周期（天），例：为1时，代表每天备份一次，为2时表示每两天备份一次。最小值为1
* start_time_hour  每次备份的时间（24小时制），例：我想凌晨2点开始备份，这里就填入2。最小值为0，最大值24
* start_time_minute  每次备份的时间（分钟），例：我想下午15点20分开始备份，这里就填20，上边的start_time_hour填15.最小值为0，最大值为60
* auto_delete  自动删除开关。打开时（true）每次备份并上传成功后后会删除本地备份的文件，否则（false）就不会自动删除
* immediate_exec  立即执行。打开时每当程序启动会立即执行一次，然后再按照计划进行
* backup_path  备份路径。这是个列表，你可以把你想要备份的路径写在这里，注意使用绝对路径


## 使用方法

确保你已经完成了如下工作：

* json配置文件已经保存好，假设`/data/c.json`
* 工作目录已经创建
* GnuPG公钥已经放好

可以直接执行二进制文件，参数如下

```shell script
backuper -config=/data/c.json
```

也可以使用systemd来管理

```editorconfig
[Unit]
Description=Backuper
After=network.target

[Service]
User=root
Group=root
ExecStart=/data/backuper -config=/data/c.json
Type=simple
Restart=on-failure
RestartSec=20

[Install]
WantedBy=multi-user.target
```

