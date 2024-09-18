package models

import (
	"time"
)

func (o *Order) CanReturned() error {
	if o.State == ReturnedState || (o.KeepUntilDate.Before(time.Now()) && o.State == AcceptState) {
		o.State = SoftDelete
	} else {
		return ErrCanReturned
	}
	return nil
}
