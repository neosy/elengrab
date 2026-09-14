package channel

import "context"

func (uc *Channel) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return uc.channelRepo().Tx(ctx, fn)
}
