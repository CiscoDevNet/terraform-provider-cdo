package main

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
)

func main() {
	client, err := client.New("https://ci.dev.lockhart.io", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJhZTk4ZDI1Zi0xMDg5LTQyODYtYTNjNS01MDVkY2I0NDMxYTIiLCJyZWFkIiwid3JpdGUiXSwiYW1yIjoicHdkIiwicm9sZXMiOlsiUk9MRV9BRE1JTiJdLCJpc3MiOiJpdGQiLCJjbHVzdGVySWQiOiIxIiwiaWQiOiI2NzJjMGE0MS1kMjAzLTQ2YzEtYmE5ZS0wNDVmYWUwYTc5ZGQiLCJzdWJqZWN0VHlwZSI6InVzZXIiLCJqdGkiOiI0MzAxNTMxMS1mYTcyLTQ1NDEtYTc4OS0yNGM2M2JmZTU1ZjYiLCJwYXJlbnRJZCI6ImFlOThkMjVmLTEwODktNDI4Ni1hM2M1LTUwNWRjYjQ0MzFhMiIsImNsaWVudF9pZCI6ImFwaS1jbGllbnQifQ.psPRQHG4UKxYxS-xEjlo40_vTnwBkEmKc-7LSoeGxjXWywFNc1cMUCtE7aENIi-HfDertAKfatmr6ZiJE-9F9Xc1etDqv7LAhFNlKtpYiVzSGPkPbfUINuDWt59Ymy3rRA25SJIuesROVx19eXjJF9IxyGMm5sYRS4H24wd50YoMRjuget_92NXeY-XjcmaL9TSGOmO-tfzMaPs2hE7IjXBcTJaI-btA8UJLczQbkmdADnLB9OFJfHArnkgDXF5hNp8JXg3rAM8UWmJrjSnClx7XLruWISaHWGbzWBE5ydGL9egxA-r2SFmoNWyPODkDRHrivL2oEVPyj46nveWjrQ")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	_, err = client.UpdateCloudFtd(ctx, cloudftd.NewUpdateInput("66827980-584b-4ba1-9ee6-a0a1692c640f", "my-new-name"))
	if err != nil {
		panic(err)
	}
}
