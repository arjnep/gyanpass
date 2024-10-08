package usecase

import (
	"fmt"
	"log"

	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/arjnep/gyanpass/internal/repository"
	"github.com/arjnep/gyanpass/pkg/notification"
	"github.com/arjnep/gyanpass/pkg/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExchangeUsecase interface {
	RequestExchange(request *entity.ExchangeRequest) (*entity.ExchangeRequest, error)
	GetExchangeRequestByID(id uuid.UUID, userID uuid.UUID) (*entity.ExchangeRequest, error)
	GetExchangeRequestsByRequestedByID(userID uuid.UUID) ([]entity.ExchangeRequest, error)
	GetExchangeRequestsByRequestedToID(userID uuid.UUID) ([]entity.ExchangeRequest, error)
	AcceptExchange(request *entity.ExchangeRequest, userID uuid.UUID) error
	DeclineExchange(request *entity.ExchangeRequest, userID uuid.UUID) error
	DeleteExchangeRequest(request *entity.ExchangeRequest, userID uuid.UUID) error
	ConfirmExchange(request *entity.ExchangeRequest, userID uuid.UUID) error
	GetExchangeRequestsByBookIDAndUserID(bookID uint, userID uuid.UUID) ([]entity.ExchangeRequest, error)
	GetExchangeRequestsByUserID(userID uuid.UUID) ([]entity.ExchangeRequest, error)
}

type exchangeUsecase struct {
	exchangeRepo        repository.ExchangeRepository
	bookRepo            repository.BookRepository
	notificationService notification.Service
}

func NewExchangeUsecase(exchangeRepo repository.ExchangeRepository, bookRepo repository.BookRepository, notificationService notification.Service) ExchangeUsecase {
	return &exchangeUsecase{exchangeRepo, bookRepo, notificationService}
}

func (u *exchangeUsecase) RequestExchange(request *entity.ExchangeRequest) (*entity.ExchangeRequest, error) {
	if u.exchangeRepo.IsSelfRequest(request.RequestedByID, request.RequestedToID) {
		return nil, response.NewBadRequestError("Cannot Request To Yourself")
	}
	canRequest, err := u.exchangeRepo.CanRequest(request.RequestedByID, request.RequestedToID)
	if err != nil {
		return nil, response.NewInternalServerError()
	}
	if !canRequest {
		return nil, response.NewConflictError("exchange request", "one request already exists with this user")
	}
	if !request.RequestedBook.IsActive {
		return nil, response.NewConflictError("book", "requested book already in exchanging process")
	}
	if !request.OfferedBook.IsActive {
		return nil, response.NewConflictError("book", "offered book already in exchanging process")
	}

	request.Status = "pending"
	request.RequestedByConfirmed = false
	request.RequestedToConfirmed = false

	err = u.exchangeRepo.Create(request)
	if err != nil {
		return nil, response.NewInternalServerError()
	}

	u.sanitizeExchangeRequest(request, request.RequestedByID)

	msg := "You have new exchange request for your book " + request.RequestedBook.Title + "."
	err = u.notificationService.SendNotification(request.RequestedToID, "exchange request", msg)
	if err != nil {
		log.Println("Failed Sending Notification of new request:", err)
		return nil, response.NewInternalServerError()
	}

	return request, nil
}

func (u *exchangeUsecase) GetExchangeRequestByID(id uuid.UUID, userID uuid.UUID) (*entity.ExchangeRequest, error) {
	request, err := u.exchangeRepo.FindByID(id)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, response.NewNotFoundError("exchange request", fmt.Sprintf("%v", id))
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, response.NewInternalServerError()
	}

	if request.RequestedByID != userID && request.RequestedToID != userID {
		// return nil, response.NewAuthorizationError("you do not have permission")
		return nil, response.NewNotFoundError("exchange request", fmt.Sprintf("%v", id))
	}

	if request.Status != "accepted" && request.Status != "exchanged" {
		u.sanitizeExchangeRequest(request, userID)
	}

	return request, nil

}

func (u *exchangeUsecase) GetExchangeRequestsByBookIDAndUserID(bookID uint, userID uuid.UUID) ([]entity.ExchangeRequest, error) {
	requests, err := u.exchangeRepo.FindRequestsByBookIDAndUserID(bookID, userID)
	if err != nil {
		return nil, response.NewInternalServerError()
	}

	for i := range requests {
		if requests[i].Status != "accepted" && requests[i].Status != "exchanged" {
			u.sanitizeExchangeRequest(&requests[i], userID)
		}
	}

	return requests, nil
}

func (u *exchangeUsecase) GetExchangeRequestsByUserID(userID uuid.UUID) ([]entity.ExchangeRequest, error) {
	requests, err := u.exchangeRepo.FindRequestsByUserID(userID)
	if err != nil {
		return nil, response.NewInternalServerError()
	}

	for i := range requests {
		if requests[i].Status != "accepted" && requests[i].Status != "exchanged" {
			u.sanitizeExchangeRequest(&requests[i], userID)
		}
	}

	return requests, nil
}

func (u *exchangeUsecase) GetExchangeRequestsByRequestedByID(userID uuid.UUID) ([]entity.ExchangeRequest, error) {
	requests, err := u.exchangeRepo.FindByRequestedByID(userID)
	if err != nil {
		return nil, response.NewInternalServerError()
	}

	for i := range requests {
		if requests[i].Status != "accepted" && requests[i].Status != "exchanged" {
			u.sanitizeExchangeRequest(&requests[i], userID)
		}
	}

	return requests, nil
}

func (u *exchangeUsecase) GetExchangeRequestsByRequestedToID(userID uuid.UUID) ([]entity.ExchangeRequest, error) {
	requests, err := u.exchangeRepo.FindByRequestedToID(userID)
	if err != nil {
		return nil, response.NewInternalServerError()
	}

	for i := range requests {
		if requests[i].Status != "accepted" && requests[i].Status != "exchanged" {
			u.sanitizeExchangeRequest(&requests[i], userID)
		}
	}

	return requests, nil
}

func (u *exchangeUsecase) AcceptExchange(request *entity.ExchangeRequest, userID uuid.UUID) error {
	if request.RequestedToID != userID {
		// return response.NewAuthorizationError("you do not have permission")
		return response.NewNotFoundError("exchange request", fmt.Sprintf("%v", request.ID))
	}
	if request.Status != "pending" {
		return response.NewBadRequestError("request already accepted or declined")
	}

	err := u.resolveExchangeRequest(request, "accepted")
	if err != nil {
		return err
	}

	return nil

}

func (u *exchangeUsecase) DeclineExchange(request *entity.ExchangeRequest, userID uuid.UUID) error {
	if request.RequestedToID != userID {
		// return response.NewAuthorizationError("you do not have permission")
		return response.NewNotFoundError("exchange request", fmt.Sprintf("%v", request.ID))

	}
	if request.Status != "pending" {
		return response.NewBadRequestError("request already accepted or declined")
	}
	err := u.resolveExchangeRequest(request, "declined")
	if err != nil {
		return err
	}

	return nil
}

func (u *exchangeUsecase) ConfirmExchange(request *entity.ExchangeRequest, userID uuid.UUID) error {
	if request.RequestedByID != userID && request.RequestedToID != userID {
		// return response.NewAuthorizationError("you do not have permission")
		return response.NewNotFoundError("exchange request", fmt.Sprintf("%v", request.ID))
	}
	if request.Status == "exchanged" {
		return response.NewBadRequestError("request is already confirmed")
	} else if request.Status == "declined" {
		return response.NewBadRequestError("request is already declined")
	} else if request.Status == "pending" {
		return response.NewBadRequestError("request is not accepted")
	}

	var recipientID uuid.UUID
	var msg string

	if request.RequestedByID == userID {
		request.RequestedByConfirmed = true
		recipientID = request.RequestedToID
		msg = request.RequestedBy.FirstName + " confirmed the exchange request."
	} else if request.RequestedToID == userID {
		request.RequestedToConfirmed = true
		recipientID = request.RequestedByID
		msg = request.RequestedTo.FirstName + " confirmed the exchange request."
	}
	err := u.resolveExchangeRequest(request, "exchanged")
	if err != nil {
		return err
	}

	err = u.notificationService.SendNotification(recipientID, "exchange request", msg)
	if err != nil {
		log.Println("Failed Sending Notification confirm:", err)
		return response.NewInternalServerError()
	}

	return nil
}

func (u *exchangeUsecase) DeleteExchangeRequest(request *entity.ExchangeRequest, userID uuid.UUID) error {
	if request.RequestedByID != userID {
		// return response.NewAuthorizationError("you do not have permission")
		return response.NewNotFoundError("exchange request", fmt.Sprintf("%v", request.ID))
	}
	if request.Status != "pending" && request.Status != "declined" {
		return response.NewBadRequestError("only pending requests can be deleted")
	}
	err := u.exchangeRepo.Delete(request)
	if err != nil {
		return response.NewInternalServerError()
	}
	if !request.RequestedBook.IsActive {
		request.RequestedBook.IsActive = true
		bookUpdates := map[string]interface{}{
			"is_active": true,
		}
		err := u.bookRepo.Update(&request.RequestedBook, bookUpdates)
		if err != nil {
			return response.NewInternalServerError()
		}
	}
	if !request.OfferedBook.IsActive {
		request.OfferedBook.IsActive = true
		bookUpdates := map[string]interface{}{
			"is_active": true,
		}
		err := u.bookRepo.Update(&request.OfferedBook, bookUpdates)
		if err != nil {
			return response.NewInternalServerError()
		}
	}

	return nil

}

func (u *exchangeUsecase) resolveExchangeRequest(request *entity.ExchangeRequest, status string) error {
	switch status {
	case "accepted":
		request.Status = "accepted"
		request.RequestedBook.IsActive = false
		request.RequestedBook.IsActive = false

		err := u.exchangeRepo.Update(request)
		if err != nil {
			return response.NewInternalServerError()
		}
		bookUpdates := map[string]interface{}{
			"is_active": false,
		}
		err = u.bookRepo.Update(&request.RequestedBook, bookUpdates)
		if err != nil {
			return response.NewInternalServerError()
		}
		err = u.bookRepo.Update(&request.OfferedBook, bookUpdates)
		if err != nil {
			return response.NewInternalServerError()
		}
		pendingRequests, err := u.exchangeRepo.FindPendingRequestsByBookID(request.RequestedBookID)
		if err != nil {
			return response.NewInternalServerError()
		}
		for _, pendingRequest := range pendingRequests {
			if pendingRequest.ID != request.ID {
				pendingRequest.Status = "declined"
				err := u.exchangeRepo.Update(&pendingRequest)
				if err != nil {
					return response.NewInternalServerError()
				} else {
					msg := "Your Exchange Request For Book " + pendingRequest.RequestedBook.Title + " is declined."
					err = u.notificationService.SendNotification(pendingRequest.RequestedByID, "exchange request", msg)
					if err != nil {
						log.Println("Failed Sending Notification for decline:", err)
						return response.NewInternalServerError()
					}
				}
			}
		}
		msg := "Your Exchange Request For Book " + request.RequestedBook.Title + " is accepted."
		err = u.notificationService.SendNotification(request.RequestedByID, "exchange request", msg)
		if err != nil {
			log.Println("Failed Sending Notification for accept:", err)
			return response.NewInternalServerError()
		}
	case "declined":
		request.Status = "declined"
		err := u.exchangeRepo.Update(request)
		if err != nil {
			return response.NewInternalServerError()
		}
		msg := "Your Exchange Request For Book " + request.RequestedBook.Title + " is declined."
		err = u.notificationService.SendNotification(request.RequestedByID, "exchange request", msg)
		if err != nil {
			log.Println("Failed Sending Notification for decline:", err)
			return response.NewInternalServerError()
		}
	case "exchanged":
		if request.RequestedByConfirmed && request.RequestedToConfirmed {
			request.Status = "exchanged"
		}
		err := u.exchangeRepo.Update(request)
		if err != nil {
			return response.NewInternalServerError()
		}
	}
	return nil
}

func (u *exchangeUsecase) sanitizeExchangeRequest(request *entity.ExchangeRequest, userID uuid.UUID) {
	request.RequestedBook.Owner.Role = ""
	request.OfferedBook.Owner.Role = ""
	if request.RequestedByID == userID {
		request.RequestedBook.PickupLocation.Latitude = 0
		request.RequestedBook.PickupLocation.Longitude = 0
		request.RequestedBook.Owner.Email = ""
		request.RequestedBook.Owner.Phone = ""
	} else if request.RequestedToID == userID {
		request.OfferedBook.PickupLocation.Latitude = 0
		request.OfferedBook.PickupLocation.Longitude = 0
		request.OfferedBook.Owner.Email = ""
		request.OfferedBook.Owner.Phone = ""
	}
}
