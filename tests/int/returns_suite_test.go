package int

import (
	"encoding/json"
	"errors"
	e "gitlab.ozon.dev/akugnerevich/homework.git/internal/service/errors"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/service/returns"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/orderStorage"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/storage/returnStorage"
)

type ReturnServiceTestSuite struct {
	suite.Suite
	orderStorage  *orderStorage.OrderStorage
	returnStorage *returnStorage.ReturnStorage
}

func (suite *ReturnServiceTestSuite) SetupTest() {
	suite.orderStorage = newOrderStorage()
	suite.returnStorage = newReturnStorage()
}

func (suite *ReturnServiceTestSuite) TestRefundOrder_Success() {
	order := &models.Order{
		ID:        1,
		UserID:    1,
		State:     models.PlaceState,
		PlaceDate: time.Now(),
	}
	suite.orderStorage.AddOrderToStorage(order)

	err := returns.RefundOrder(suite.returnStorage, suite.orderStorage, 1, 1)
	suite.NoError(err)

	suite.Equal(models.ReturnedState, order.State)
	suite.True(suite.returnStorage.IsConsist(1))
}

func (suite *ReturnServiceTestSuite) TestRefundOrder_NoOrder() {
	err := returns.RefundOrder(suite.returnStorage, suite.orderStorage, 1, 1)
	suite.True(errors.Is(err, e.ErrCheckOrderID))
}

func (suite *ReturnServiceTestSuite) TestListReturn() {
	suite.returnStorage.AddReturnToStorage(&models.Return{ID: 1, UserID: 1})

	err := returns.ListReturns(suite.returnStorage, 10, 1)
	suite.NoError(err)
}

func (suite *ReturnServiceTestSuite) TestWriteToJSON() {
	returnItem := &models.Return{ID: 1, UserID: 1}
	suite.returnStorage.AddReturnToStorage(returnItem)

	err := suite.returnStorage.WriteToJSON()
	suite.NoError(err)

	file, err := os.Open(suite.returnStorage.GetPath())
	suite.NoError(err)
	defer file.Close()
	defer os.Remove(suite.returnStorage.GetPath())

	var storageData returnStorage.ReturnStorage
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&storageData)
	suite.NoError(err)

	storedReturn, exists := storageData.GetReturn(returnItem.ID)
	suite.True(exists)
	suite.Equal(returnItem.UserID, storedReturn.UserID)
}

func (suite *ReturnServiceTestSuite) TestReadFromJSON() {
	returnItem := &models.Return{ID: 1, UserID: 1}
	suite.returnStorage.AddReturnToStorage(returnItem)

	err := suite.returnStorage.WriteToJSON()
	suite.NoError(err)

	suite.returnStorage.Data = make(map[uint]*models.Return)

	err = suite.returnStorage.ReadFromJSON()
	suite.NoError(err)

	defer os.Remove(suite.returnStorage.GetPath())

	storedReturn, exists := suite.returnStorage.GetReturn(returnItem.ID)
	suite.True(exists)
	suite.Equal(returnItem.UserID, storedReturn.UserID)
}

func TestReturnServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReturnServiceTestSuite))
}
