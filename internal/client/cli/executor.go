package cli

import (
	"fmt"
	"github.com/mdp/qrterminal/v3"
	uuid "github.com/satori/go.uuid"
	"os"
	"strings"

	"github.com/unbeman/ya-prac-go-second-grade/internal/client"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/model"
)

type exec struct {
	app      *client.ClientApp
	commands map[string]func(input []string)
}

func GetExecutor(app *client.ClientApp) *exec {
	e := &exec{}
	e.app = app
	commands := map[string]func(agrs []string){
		"exit":         e.exit,
		"register":     e.register,
		"login":        e.login,
		"genTOTP":      e.generate2FA,
		"verifyTOTP":   e.verify2FA,
		"validateTOTP": e.validateTOTP,
		"disableTOTP":  e.disableTOTP,
		"types":        e.types,
		"create":       e.create,
		"edit":         e.edit,
		"get-all":      e.getAll,
		"get":          e.get,
		"search":       e.search,
		"delete":       e.delete,
	}
	e.commands = commands

	return e
}

func (e *exec) Execute(s string) {
	args := getArgs(s)
	command := args[0]

	if v, ok := e.commands[command]; ok {
		v(args)
	}

}

func (e exec) exit(args []string) {
	fmt.Println("goodbye!")
	e.app.Stop()
	os.Exit(0)
}

func (e exec) register(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: Login and Password is missing")
		return
	case 1:
		fmt.Println("validation error: Password is missing")
		return
	}

	err := e.app.Auth.Register(args[1], args[2])
	if err != nil {
		fmt.Println("error occurred: ", err)
	}
}

func (e exec) login(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: Login and Password is missing")
		return
	case 1:
		fmt.Println("validation error: Password is missing")
		return
	}

	err := e.app.Auth.Login(args[1], args[2])
	if err != nil {
		fmt.Println("error occurred: ", err)
		return
	}
	fmt.Println("please validate 2fa code")
}

func (e exec) generate2FA(args []string) {
	secret, url, err := e.app.F2a.Generate()
	if err != nil {
		fmt.Println("error occurred: ", err) //todo: you should login first
		return
	}

	qrterminal.Generate(url, qrterminal.L, os.Stdout)
	fmt.Println()
	qrterminal.GenerateHalfBlock(url, qrterminal.L, os.Stdout)

	fmt.Println("Secret key: ", secret)
	fmt.Println("URL: ", url)
}

func (e *exec) verify2FA(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: TOTP token required")
		return
	}
	err := e.app.F2a.Verify(args[1])
	if err != nil {
		fmt.Println("error occurred: ", err)
	} else {
		fmt.Println("successful verified")
	}
	e.app.Sync.StartSync()
}

func (e *exec) validateTOTP(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: TOTP token required")
		return
	}
	err := e.app.F2a.Validate(args[1])
	if err != nil {
		fmt.Println("error occurred: ", err)
	} else {
		fmt.Println("successful validated")
	}

	e.app.Sync.StartSync()
}

func (e *exec) disableTOTP(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: login required")
		return
	case 1:
		fmt.Println("validation error: password required")
		return
	}
	err := e.app.F2a.Disable(args[1], args[2])
	if err != nil {
		fmt.Println("error occurred: ", err)
	} else {
		fmt.Println("successful disabled 2FA")
	}
}

func (e exec) types(args []string) {
	types := e.app.Sync.GetTypes()
	for t, _ := range types {
		fmt.Println(t)
	}
}

func (e exec) create(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: type required")
		return
	}
	choosedType := model.CredentialType(args[1])

	var meta, sensitive string
	var err error

	switch choosedType {
	case model.Login:
		meta, sensitive, err = inputLogin(args)
	case model.Bank:
		meta, sensitive, err = inputBank(args)
	case model.Note:
		meta, sensitive, err = inputNote(args)
	default:
		fmt.Println("unsupported type")
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	err = e.app.Sync.CreateCred(choosedType, meta, sensitive)
	if err != nil {
		fmt.Println("error occurred: ", err)
		return
	} else {
		fmt.Println("successfully created")
	}

}

func (e exec) edit(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: id required")
		return
	}

	credID, err := uuid.FromString(args[1])
	if err != nil {
		fmt.Println("validation error: not uuid")
		return
	}
	cred, err := e.app.Sync.GetCredByID(credID)
	if err != nil {
		fmt.Println("error occurred: ", err)
		return
	}
	var meta, sensitive string
	switch cred.Type {
	case model.Login:
		meta, sensitive, err = inputLogin(args)
	case model.Bank:
		meta, sensitive, err = inputBank(args)
	case model.Note:
		meta, sensitive, err = inputNote(args)
	default:
		fmt.Println("unsupported type")
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	err = e.app.Sync.EditCred(cred, meta, sensitive)
	if err != nil {
		fmt.Println("error occurred: ", err)
		return
	} else {
		fmt.Println("successfully updated")
	}

}

func (e exec) getAll(args []string) {
	creds, err := e.app.Sync.GetAll()
	if err != nil {
		fmt.Println("error occurred: ", err)
	}

	fmt.Println("ID 					Type 	Info")

	for _, cred := range creds {
		switch cred.Type {
		case model.Login:
			printLoginMeta(cred)
		case model.Bank:
			printBankMeta(cred)
		case model.Note:
			printNoteMeta(cred)
		}
	}

	fmt.Println("Get row by id to show secret")
}

func (e exec) get(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: index required")
		return
	}

	credID, err := uuid.FromString(args[1])
	if err != nil {
		fmt.Println("validation error: not id")
		return
	}

	cred, err := e.app.Sync.GetCredByID(credID)
	if err != nil {
		fmt.Println("error occurred: ", err)
		return
	}

	switch cred.Type {
	case model.Login:
		printLogin(cred)
	case model.Bank:
		printBank(cred)
	case model.Note:
		printNote(cred)
	}
}

func (e exec) search(args []string) {
	switch len(args) - 1 {
	case 0:
		fmt.Println("validation error: search query required")
		return
	}

	creds, err := e.app.Sync.Search(args[1])
	if err != nil {
		fmt.Println("error occurred: ", err)
		return
	}

	fmt.Printf("found %d records \n", len(creds))

	fmt.Println("ID 					Type 	Info")

	for _, cred := range creds {
		switch cred.Type {
		case model.Login:
			printLoginMeta(cred)
		case model.Bank:
			printBankMeta(cred)
		case model.Note:
			printNoteMeta(cred)
		}
	}
	fmt.Println("Get row by id to show secret")
}

func (e exec) delete(args []string) {

}

func getArgs(s string) []string {
	s = strings.TrimSpace(s)
	return strings.Split(s, " ")
}

func inputLogin(args []string) (string, string, error) {
	var meta, sensitive string
	switch len(args) - 2 {
	case 0:
		return meta, sensitive, fmt.Errorf("validation error: site required")
	case 1:
		return meta, sensitive, fmt.Errorf("validation error: login required")
	case 2:
		return meta, sensitive, fmt.Errorf("validation error: password required")
	}

	site := args[2]
	login := args[3]
	password := args[4]

	meta = fmt.Sprintf("%s:%s", site, login)
	sensitive = password
	return meta, sensitive, nil
}

func inputBank(args []string) (string, string, error) {
	var meta, sensitive string
	switch len(args) - 2 {
	case 0:
		return meta, sensitive, fmt.Errorf("validation error: bank name required")
	case 1:
		return meta, sensitive, fmt.Errorf("validation error: credit card number required")
	case 2:
		return meta, sensitive, fmt.Errorf("validation error: expired at required")
	case 3:
		return meta, sensitive, fmt.Errorf("validation error: cvc required")
	}
	info := args[2]
	number := args[3]
	expired := args[4]
	cvc := args[5]

	meta = info
	sensitive = fmt.Sprintf("%s:%s:%s", number, expired, cvc)
	return meta, sensitive, nil
}

func inputNote(args []string) (string, string, error) {
	var meta, sensitive string
	switch len(args) - 2 {
	case 0:
		return meta, sensitive, fmt.Errorf("validation error: note title required")
	case 1:
		return meta, sensitive, fmt.Errorf("validation error: text (in one line) required")
	}
	title := args[2]
	note := args[3:]
	meta = title

	sensitive = strings.Join(note, " ")
	return meta, sensitive, nil
}

func printLoginMeta(cred model.Credential) {
	splittedMeta := strings.Split(cred.MetaData, ":")

	site, login := splittedMeta[0], splittedMeta[1]
	fmt.Printf("%s	login	Site: %s	Login: %s\n", cred.ID, site, login)
}

func printBankMeta(cred model.Credential) {
	fmt.Printf("%s	bank	Card Name: %s\n", cred.ID, cred.MetaData)
}

func printNoteMeta(cred model.Credential) {
	fmt.Printf("%s	note	Title: %s\n", cred.ID, cred.MetaData)
}

func printLogin(cred model.Credential) {
	splittedMeta := strings.Split(cred.MetaData, ":")

	site, login := splittedMeta[0], splittedMeta[1]
	fmt.Printf("Site: %s\n	Login: %s\n		Password: %s\n", site, login, string(cred.Decrypted))
}

func printBank(cred model.Credential) {
	raw := string(cred.Decrypted)
	data := strings.Split(raw, ":")
	cardNumber, exp, cvc := data[0], data[1], data[2]
	fmt.Printf("Card Name: %s\nNumber: %s\nExp: %s\nCVC: %s\n", cred.MetaData, cardNumber, exp, cvc)
}

func printNote(cred model.Credential) {
	raw := string(cred.Decrypted)
	fmt.Printf("Title: %s\nNote:\n%s\n", cred.MetaData, raw)
}
