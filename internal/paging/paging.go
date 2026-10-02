package paging

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type PageParameters struct {
	size   int
	number int
}

func (pp *PageParameters) ConfigureFlags(fs *pflag.FlagSet) {
	fs.IntVar(&pp.size, "limit", 100, "Number of entries to receive at most.")
	fs.IntVar(&pp.number, "page", 0, "Page number to calculate first item to receive. Page numbers start from `1`.")

	commands.Must(fs.SetAnnotation("limit", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("page", commands.FlagAnnotationNoFileCompletions, nil))
}

func (pp *PageParameters) Page() *request.Page {
	return &request.Page{
		Number: pp.number,
		Size:   pp.size,
	}
}

// LimitOffset returns the page as limit and offset query values, computed the same way as the v8 SDK.
func (pp *PageParameters) LimitOffset() (limit, offset int) {
	limit = max(pp.size, 0)
	return limit, (max(pp.number, 1) - 1) * limit
}
