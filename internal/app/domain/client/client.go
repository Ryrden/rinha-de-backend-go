package client

type Client struct {
	ID      string
	Limit   int
	Balance int
}

func (c *Client) AddTransaction(value int, kind string) {
	switch kind {
	case "c":
		c.Balance += value
	case "d":
		c.Balance -= value
	}
}

func (c *Client) CanAfford(value int, kind string) bool {
	if kind == "d" && c.Balance-value < -c.Limit {
		return false
	}
	return true
}

func NewClient(id string, limit int, balance int) *Client {
	return &Client{
		ID:      id,
		Limit:   limit,
		Balance: balance,
	}
}
