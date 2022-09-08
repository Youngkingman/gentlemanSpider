# gentlemanSpider
*Get some Hons my friends！GKDGKD！*

## 用途：

- 爬取并保存网站 [绅士漫画](https://www.wnacg.com) 的本子/色图（大陆用户要求在梯子的环境下使用）
- 支持按照配置进行多线程(协程)加速
- 按照配置对个人喜好的xp标签进行过滤（可用标签在项目下的`activeTags`文件夹中可查看）

## 项目构建：

### 直接构建（推荐）

在`go1.16`以上的环境下，中文用户参考 [Go下载 ](https://studygolang.com/dl)执行如下操作运行爬虫程序：

```bash
git clone https://github.com/Youngkingman/gentlemanSpider.git
go build
./gentlemanSpider
```

### 使用docker进行构建（不推荐）

// TODO

### 使用`go`相关工具

// TODO

### 直接下载发布版本（推荐）

// TODO

## 项目配置使用说明：

项目需要在可执行文件所在同级文件夹下新建一个`config.yaml`文件，可以直接克隆本仓库进行使用，发布版本压缩包中自带该文件。

```yaml
CrawlerSetting:
  PageStart: 2
  PageEnd: 2
  EnableProxy: true
  ProxyHost: http://127.0.0.1:62340
  TagConsumerCount: 1
  HonConsumerCount: 2
  HonBuffer: 8192
  TagBuffer: 8192
  EnableFilter: false
  WantedTags:
    - 百合
    - 女同士
```

- `PageStart`: 最小为 1 爬虫开始页面，浏览 [这里](https://www.wnacg.com/albums-index-page-1.html)查看最大页数
- `PageEnd`: 最小为 1， 必须**大于** `PageStart`，爬虫结束页面，同上浏览 [这里](https://www.wnacg.com/albums-index-page-1.html)查看最大页数
- `EnableProxy`: 是否使用代理（使用需要填写`ProxyHost`)，大陆用户填`true`，海外用户填`false`
- `ProxyHost`: 查看自己电脑的代理，需要带上端口
- `TagConsumerCount`: 用于收集本子标签的配置，配置完成后会在可执行文件同级目录下创建包含所有标签(`xp`/作者/是否彩页汉化)信息的`activeTag`文件，会一定程度减慢本子的下载，不需要可以设置为0
- `HonConsumerCount`:必须大于 0，根据程序运行电脑的配置进行设置
- `HonBuffer`:必须大于 0，
- `TagBuffer`:若设置了 `TagConsumerCount`则 必须大于 0
- `EnableFilter`:是否使用标签过滤，`true`为使用，`false`会默认不过滤爬取所有本子
- `WantedTags`: 你的标签集，需要`EnableFilter`为`true`，每行以`-{yourXP}`的形式对标签进行选取，可以在`activeTag`文件中选取自己所需的标签，最终下载的本子至少会具有你给出的标签集中的一个标签（标签给的越少下载的本子一般会越少）

## 如何查看自己电脑的代理服务器端口(Win)

Linus用户请自行探索翻墙设置。

// TODO
