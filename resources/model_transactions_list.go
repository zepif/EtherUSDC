/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type TransactionsList struct {
	Key
}
type TransactionsListResponse struct {
	Data     TransactionsList `json:"data"`
	Included Included         `json:"included"`
}

type TransactionsListListResponse struct {
	Data     []TransactionsList `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
}

// MustTransactionsList - returns TransactionsList from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustTransactionsList(key Key) *TransactionsList {
	var transactionsList TransactionsList
	if c.tryFindEntry(key, &transactionsList) {
		return &transactionsList
	}
	return nil
}
