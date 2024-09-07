package models

import (
	"fmt"
	"time"
)

func (o *Order) CanReturned() error {
	if o.State == ReturnedState || (o.KeepUntilDate.Before(time.Now()) && o.State == AcceptState) {
		o.State = SoftDelete
	} else {
		return fmt.Errorf("Order can`t be returned because it has not expired yet and this order is not a return \n")
	}
	return nil
}
