package libs

type Options struct {
	Logfile     string   //日志文件路径
	Target      string   //需要单个扫描的目标
	TargetFile  string   //需要扫描的文件列表
	Urls        []string //需要扫描的列表-处理后的
	Threads     int      //单个任务的扫描并发线程数
	Concurrency int      //并行任务数
	Timeout     int      //超时时间
	FilterCode  string   //过滤的状态码
	Wordlist    string   //字典路径
	Method      string   //请求方法
	OutPutFile  string   //输出结果文件
}
