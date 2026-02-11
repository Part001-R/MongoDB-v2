package mongodb

import (
	"context"
	"fmt"
	"time"
)

// Close db connect. Return error.
func (d *mongoDB) Close() error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := d.db.Disconnect(ctx); err != nil {
		return fmt.Errorf("Function Disconnetc, return error <%w>", err)
	}

	return nil
}
