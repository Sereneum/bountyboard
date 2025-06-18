package task

type Task struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Done         bool   `json:"done"`
	BountyAmount int    `json:"bounty_amount"`
	UserID       string `json:"user_id"`
}
