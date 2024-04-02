package paging

import (
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type PageParameters struct {
	size   int
	number int
}

func (pp *PageParameters) ConfigureFlags(fs *pflag.FlagSet) {
	fs.IntVar(&pp.size, "size", 100, "Number of entries to receive at most.")
	fs.IntVar(&pp.number, "page", 0, "Page number to calculate first item to receive. Page numbers start from `1`.")
}

func (pp *PageParameters) Page() *request.Page {
	return &request.Page{
		Number: pp.number,
		Size:   pp.size,
	}
}
