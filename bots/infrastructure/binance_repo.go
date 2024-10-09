package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/juankohler/crypto-bot/bots/domain"
	"github.com/juankohler/crypto-bot/libs/go/errors"
	"github.com/juankohler/crypto-bot/libs/go/logs"
	"github.com/juankohler/crypto-bot/libs/go/models"
	"github.com/juankohler/crypto-bot/libs/go/restclient"
	"github.com/shopspring/decimal"
)

type repository struct {
	getPriceEndpoint restclient.Endpoint
}

func NewBinanceRepo(config *restclient.Config) (*repository, error) {
	client := restclient.New(*config)

	failAtInternalErrorCodes := func(req restclient.Request, res restclient.Response) error {
		if res.StatusCode() < 200 || res.StatusCode() >= 300 {
			bodyString := string(res.Body())
			bodyString = strings.ReplaceAll(bodyString, "\n", "")
			bodyString = strings.ReplaceAll(bodyString, "\r", "")
			bodyString = strings.ReplaceAll(bodyString, "  ", " ")

			return errors.New(
				domain.ErrInternal,
				fmt.Sprintf("Request failed with status code %d. body: %s", res.StatusCode(), bodyString),
			)
		}
		return nil
	}

	repo := &repository{
		getPriceEndpoint: client.GET(
			"/v3/ticker/price",
			restclient.Header("content-type", "application/json"),
			restclient.FailAt(failAtInternalErrorCodes),
		),
	}

	return repo, nil
}

type GetPriceResponse struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
}

func (r *repository) GetPrice(ctx context.Context, baseCurrency string, quoteCurrency string) (*domain.Price, error) {
	/** Build EndpointOptions */
	var restclientOptions []restclient.EndpointOption
	restclientOptions = append(
		restclientOptions,
		restclient.QueryParam("symbol", baseCurrency+quoteCurrency),
	)

	/** Do request */
	res := r.getPriceEndpoint.DoRequest(
		ctx,
		restclientOptions...,
	)
	if res.Err() != nil {
		if res.StatusCode() == 404 {
			return nil, errors.Wrap(domain.ErrNotFound, res.Err(), "Failed to do request.")
		}

		return nil, errors.Wrap(domain.ErrInternal, res.Err(), "Failed to do request.")
	}

	var respMsg GetPriceResponse
	err := json.Unmarshal(res.Body(), &respMsg)
	if err != nil {
		return nil, errors.Wrap(domain.ErrInternal, err, fmt.Sprintf("Failed to unmarshal item. body: %s", string(res.Body())))
	}

	entity, err := domain.NewPrice(
		baseCurrency,
		quoteCurrency,
		respMsg.Price,
	)
	if err != nil {
		return nil, errors.Wrap(domain.ErrInternal, err, fmt.Sprintf("failed to parse to entity. body: %s", string(res.Body())))
	}

	return entity, nil
}

func (r *repository) CreateOrderInProvider(ctx context.Context, order *domain.Order, botName string) (string, error) {
	logs.Info(ctx, fmt.Sprintf("%s: Compra %s BTC a %s USDT (%s USDT), take_profit_price: %s", botName, order.Quantity.String(), order.EntryPrice.String(), order.InitialQuoteAmount.String(), order.TakeProfitPrice.String()))
	return models.GenerateUUID().String(), nil
}
