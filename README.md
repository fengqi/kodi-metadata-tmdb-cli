# kodi-metadata-tmdb-cli

Kodi 刮削器命令行版本，使用TMDB数据源匹配剧集和电影并生成Kodi兼容的nfo文件。
Kodi需要配置为仅本地，开发版本为Kodi 19，其他版本自行测试。

# 注意事项

1. 多个搜索结果
如搜索：Fall In Love 会出现多条，无法很好的确定是哪一个，
可通过在剧集的根目录创建 `tmdb-id.txt` 并填入在 `tmdb` 上找到的剧集 id 来指定。

2. 分季命名不规范
如 `进击的巨人` 有部分站将第二的文件名接着上一季，而不是从 `S2E01` 开始会导致匹配失败。

3. 接口调用失败
表现为所有剧集都不能匹配，大概率是 `https://api.themoviedb.org/` 域名被阻断，请自行科学上网或修改HOSTS解决。

4. 运行报错 `Error No more hard-disk space available` 或 `no space left on device` 但实际上磁盘空间足够
其实原因是inotify的允许监听文件数不够用了，可通过增大 `fs.inotify.max_user_watches` 解决。


# 更新线路图

- [x] 建议日志打印和分级
- [x] 从TMDB获取电视剧信息
- [x] 从TMDB获取电视剧分集信息
- [x] 从TMDB获取电视剧演员列表信息
- [x] 从TMDB下载封面、还报等图片
- [x] 定时扫描电视剧
- [x] 实时监听新添加的电视剧
- [x] 命名不规范的电视剧支持指定tv_id
- [x] 命名不规范的电视剧支持指定season
- [ ] 命名不规范的电视剧支持指定season和episode的对应，如：E21-34其实是第2季
- [ ] 从TMDB获取电视剧内容分级信息
- [x] 定时扫描电影
- [x] 实时监听新添加的电影
- [ ] 多个搜索结果根据集数尝试确定
- [ ] 适配其他数据源，如：imdb、tvdb，以补全部分tmdb没有的数据
- [ ] 支持嵌套匹配，如电视剧、电影合集
- [ ] 更新nfo后触发kodi更新数据
- [x] 支持单个电影文件和目录
- [x] 识别蓝光电影目录
- [x] 支持 .part 和 .!qb 文件

# 参考

- Kodi v19 (Matrix) JSON-RPC API/V12 https://kodi.wiki/view/JSON-RPC_API/v12
- Kodi v19 (Matrix) NFO files https://kodi.wiki/view/NFO_files
- TMDB Api Overview https://www.themoviedb.org/documentation/api
- TMDB Api V3 https://developers.themoviedb.org/3/getting-started/introduction
- File system notifications for Go https://github.com/fsnotify/fsnotify
