/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type TransactionResponse struct {
	Key
}
type TransactionResponseResponse struct {
	Data     TransactionResponse `json:"data"`
	Included Included            `json:"included"`
}

type TransactionResponseListResponse struct {
	Data     []TransactionResponse `json:"data"`
	Included Included              `json:"included"`
	Links    *Links                `json:"links"`
}

// MustTransactionResponse - returns TransactionResponse from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustTransactionResponse(key Key) *TransactionResponse {
	var transactionResponse TransactionResponse
	if c.tryFindEntry(key, &transactionResponse) {
		return &transactionResponse
	}
	return nil
}
