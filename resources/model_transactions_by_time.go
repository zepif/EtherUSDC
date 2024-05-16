/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type TransactionsByTime struct {
	Key
}
type TransactionsByTimeResponse struct {
	Data     TransactionsByTime `json:"data"`
	Included Included           `json:"included"`
}

type TransactionsByTimeListResponse struct {
	Data     []TransactionsByTime `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
}

// MustTransactionsByTime - returns TransactionsByTime from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustTransactionsByTime(key Key) *TransactionsByTime {
	var transactionsByTime TransactionsByTime
	if c.tryFindEntry(key, &transactionsByTime) {
		return &transactionsByTime
	}
	return nil
}
