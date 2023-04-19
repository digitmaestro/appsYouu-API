package mailchimp

import (
	"github.com/hanzoai/gochimp3"
)

/* This string must be at most 10 characters long. */
const UserIdFieldId = "KES_USERID"

func AssertMergeFields(apiKey string, listId string) {
	client := gochimp3.New(apiKey)

	list, err := client.GetList(listId, nil)

	if err != nil {
		return
	}

	mergeFields, err := list.GetMergeFields(nil)

	if err != nil {
		return
	}

	found := false

	for _, mergeField := range mergeFields.MergeFields {
		found = found || mergeField.Tag == UserIdFieldId
	}

	if !found {
		list.CreateMergeField(&gochimp3.MergeFieldRequest{
			Name:     "Keskin User ID",
			Tag:      UserIdFieldId,
			Type:     "text",
			Public:   false,
			Required: false,
		})
	}
}
