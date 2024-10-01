package models

import (
	"time"
)

func (o *Order) CanBeReturned() error {
	if o.State == RefundedState || (o.KeepUntilDate.Before(time.Now()) && o.State == AcceptState) {
		return nil
	} else {
		return ErrCanReturned
	}

}
