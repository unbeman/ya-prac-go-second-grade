package cli

import "github.com/c-bata/go-prompt"

func Complete(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest

	if d.FindStartOfPreviousWord() == 0 {
		s = []prompt.Suggest{
			{Text: "register", Description: "Register new user"},
			{Text: "login", Description: "Authenticate user"},
			{Text: "genTOTP", Description: "Get 2FA setup"},
			{Text: "verifyTOTP", Description: "Verify 2FA setup"},
			{Text: "validateTOTP", Description: "Validate OTP code"},
			{Text: "disableTOTP", Description: "Disable 2FA feature"},
			{Text: "types", Description: "Get list of credential types"},
			{Text: "create", Description: "Create new credential"},
			{Text: "edit", Description: "Update credential"},
			{Text: "get", Description: "Get by row number"},
			{Text: "get-all", Description: "Get secrets list"},
			{Text: "search", Description: "Search credential by metadata"},
			{Text: "delete", Description: "Delete credential"},
			{Text: "get-by-type", Description: "Get list of credentials by type"},
			{Text: "exit", Description: "Exit app"},
		}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
