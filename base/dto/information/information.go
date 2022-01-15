package information

type Information struct {
	TotalUserNum 		int64	`json:"totalUserNum"`
	TotalArticleNum 	int64	`json:"totalArticleNum"`
	TotalChatPersonNum	int64	`json:"totalChatPersonNum"`
	TotalReadNum		int64	`json:"totalReadNum"`
}

const (
	TOTAL_USER_NUM = "totalUserNum"
	TOTAL_ARTICLE_NUM = "totalArticleNum"
	TOTAL_READ_NUM = "totalReadNum"
)