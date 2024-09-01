package utils

import "fmt"

func Generate(account_id, formId string) string {
	return fmt.Sprintf("http://%s.localhost/to/%s", account_id, formId)
}
