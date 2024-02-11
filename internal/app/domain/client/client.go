package client

type Client struct {
	ID      string
	Limit   int
	Balance int
}

func (c *Client) AddTransaction(value int) {
	c.Balance += value
}

func (c *Client) SubtractTransaction(value int) {
	c.Balance -= value
}

func (c *Client) CanAfford(value int) bool {
	return c.Balance >= value
}

func (c *Client) CanAffordWithLimit(value int) bool {
	return c.Balance+value <= c.Limit
}

func (c *Client) GetBalance() int {
	return c.Balance
}

func (c *Client) GetLimit() int {
	return c.Limit
}

func (c *Client) GetID() string {
	return c.ID
}

func NewClient(id string, limit int, balance int) *Client {
	return &Client{
		ID:      id,
		Limit:   limit,
		Balance: balance,
	}
}
