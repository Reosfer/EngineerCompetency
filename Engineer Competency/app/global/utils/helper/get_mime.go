package helper

func GetMimeExcel(fileExtension string) string {
	var mime string

	if fileExtension == "xls" || fileExtension == "csv" || fileExtension == "xlsx" || fileExtension == "xlx" {
		mime = "vnd.ms-excel"
	} else {
		mime = "application"
	}

	return mime
}
