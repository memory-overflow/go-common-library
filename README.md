- [go common library](#go-common-library)
  - [httpcall 模块](#httpcall-模块)
  - [misc 模块](#misc-模块)
  - [stl\_extension 模块](#stl_extension-模块)
  - [text 模块](#text-模块)
  - [任务调度模块](#任务调度模块)

# go common library
本仓库整理了 golang 中最常用的代码、函数、和模块。

## httpcall 模块
常用的简单 http 调用的封装。
详见 [httpcall 模块说明](https://github.com/memory-overflow/go-common-library/blob/main/httpcall/readme.md)。

## misc 模块
常用的公共函数库，主要包含时间、协程安全、重试、id生成，log 相关的常用的处理函数。
详见[misc 模块说明](https://github.com/memory-overflow/go-common-library/blob/main/misc/readme.md)。

## stl_extension 模块
针对 golang 现有的 stl 的不足的扩展。主要是对 WaitGroup 和 map 的扩展。
- LimitWaitGroup -- 对于系统 WaitGroup 的扩展，支持 limit 并发限制并且阻塞。
- OrderedMap -- 实现了 c++ 中的排序 map，可以按照顺序遍历所有元素。


详见[stl_extension 模块说明](https://github.com/memory-overflow/go-common-library/blob/main/stl_extension/readme.md)。

## text 模块
常用文本处理相关方法。
- AcTrie：ac 自动机，多模式串快速匹配。在一个文本中快速找出来出现过哪些字符串子串以及其定位。可以理解对同一文本 s 多次调用 strings.Contains(s, xxx) 的加速。
- Levenshtein：计算文本编辑距离。
- TextSim：计算两个文本的相似度。


详见[text 模块说明](https://github.com/memory-overflow/go-common-library/blob/main/text/readme.md)。


## 任务调度模块
任务调度模块在 go 语言中，提供了一个轻量级的任务调度框架，方便对各种同步、异步任务做统一的任务调度。致力于提高构建一套任务管理调度系统的效率。
详见[任务调度框架设计原理和使用说明](https://github.com/memory-overflow/go-common-library/blob/main/task_scheduler/readme.md) 。