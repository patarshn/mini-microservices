package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"transaction-service/internal/constant"
)

type IProductRepository interface {
	GetProductBySKU(ctx context.Context, sku string) (*Product, error)
}

type productService struct {
	baseURL string
	client  *http.Client
}

func NewProductRepository(baseURL string, client *http.Client) IProductRepository {
	return &productService{
		baseURL: baseURL,
		client:  client,
	}
}

func (p *productService) GetProductBySKU(ctx context.Context, sku string) (*Product, error) {

	url := fmt.Sprintf("%s/products/sku/%s", p.baseURL, sku)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if bearerToken, ok := ctx.Value(constant.BEARER_TOKEN).(string); !ok {
		fmt.Println(bearerToken)
		return nil, errors.New("Token not in ctx")
	} else {
		req.Header.Set("Authorization", bearerToken)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseBody ResponseBody
	responseBody.Data = &Product{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	if responseBody.Error {
		return nil, fmt.Errorf("product service error: %s", responseBody.Message)
	}

	if responseBody.Data == nil {
		return nil, fmt.Errorf("product service error: data is nil")
	}

	product := responseBody.Data.(*Product)

	return product, nil
}
