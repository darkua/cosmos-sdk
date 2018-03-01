package keys

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	keys "github.com/tendermint/go-crypto/keys"

	"github.com/spf13/cobra"
)

func updateKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Change the password used to protect private key",
		RunE:  runUpdateCmd,
	}
	return cmd
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a name for the key")
	}
	name := args[0]

	buf := client.BufferStdin()
	oldpass, err := client.GetPassword(
		"Enter the current passphrase:", buf)
	if err != nil {
		return err
	}
	newpass, err := client.GetCheckPassword(
		"Enter the new passphrase:",
		"Repeat the new passphrase:", buf)
	if err != nil {
		return err
	}

	kb, err := GetKeyBase()
	if err != nil {
		return err
	}
	err = kb.Update(name, oldpass, newpass)
	if err != nil {
		return err
	}
	fmt.Println("Password successfully updated!")
	return nil
}

// REST

type UpdateKeyBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

func UpdateKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var kb keys.Keybase
	var m UpdateKeyBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	kb, err = GetKeyBase()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// TODO check if account exists and if password is correct
	err = kb.Update(name, m.OldPassword, m.NewPassword)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
}