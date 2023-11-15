package sec

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type ReadByNameInput struct {
	Name string
}

type ReadByNameOutput = ReadOutput

func ReadByName(ctx context.Context, client http.Client, inp ReadByNameInput) (*ReadByNameOutput, error) {

	client.Logger.Println("reading SEC by name")

	outp, err := ReadAll(ctx, client, ReadAllInput{})
	if err != nil {
		return nil, err
	}

	for _, sec := range *outp {
		client.Logger.Printf("checking if sec.Name %q == inp.Name %q\n", sec.Name, inp.Name)
		if sec.Name == inp.Name {
			return &sec, nil
		}
	}

	return nil, fmt.Errorf("SEC with name: %q not found", inp.Name)
}
