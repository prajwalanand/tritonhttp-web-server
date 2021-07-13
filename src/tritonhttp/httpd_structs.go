package tritonhttp

type HttpServer	struct {
	ServerPort	string
	DocRoot		string
	MIMEPath	string
	MIMEMap		map[string]string
}

type HttpResponseHeader struct {
	InitialLine string
	URL string
	FieldMap map[string]string
	ResponseBody string `default:""`
}

type HttpRequestHeader struct {
	//InitialLine string
	URL string
	FieldMap map[string]string
}