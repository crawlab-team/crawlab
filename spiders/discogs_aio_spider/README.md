```
爬取策略：
discogs抓取范围2000-2003年的
因为discogs限制显示最大10000条，所以这里采取以下策略。
1。首先抓对应的年代结果页
2。获取左栏对应的url和总数构建url
url构成
format_exact：类型
layout：sm
country_exact：国家
style_exact：风格
limit:250
year:2000-2003
decade:2000
```
### url构成
1.发现左侧区域 我们需要是style，format和country
2.根据每个分类的个数 除以每页最大显示数250，确定翻几页。