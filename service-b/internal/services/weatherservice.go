package services

import (
	"context"
	"fmt"
	"regexp"

	"github.com/teichx/go-expert-weather/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var ErrInvalidZipcode = fmt.Errorf("invalid zipcode")
var ErrCanNotFindZipcode = fmt.Errorf("can not find zipcode")

type ViaCEPClient interface {
	GetLocation(zipcode string) (string, error)
}

type WeatherAPIClient interface {
	GetTemperature(location string) (float64, error)
}

type WeatherService struct {
	viaCEP     ViaCEPClient
	weatherAPI WeatherAPIClient
}

func New(viaCEP ViaCEPClient, weatherAPI WeatherAPIClient) *WeatherService {
	return &WeatherService{
		viaCEP:     viaCEP,
		weatherAPI: weatherAPI,
	}
}

func (s *WeatherService) GetWeatherByZipcode(ctx context.Context, zipcode string) (*models.Temperature, error) {
	if err := validateZipcode(zipcode); err != nil {
		return nil, err
	}

	tracer := otel.Tracer("service-b")

	ctx, span := tracer.Start(ctx, "fetch-cep")
	span.SetAttributes(attribute.String("zipcode", zipcode))
	city, err := s.viaCEP.GetLocation(zipcode)
	span.SetAttributes(attribute.String("city", city))

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, ErrCanNotFindZipcode
	}
	span.End()

	ctx, span = tracer.Start(ctx, "fetch-weather")
	span.SetAttributes(attribute.String("city", city))
	tempC, err := s.weatherAPI.GetTemperature(city)
	span.SetAttributes(attribute.Float64("temp_C", tempC))

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, fmt.Errorf("can not find weather data: %w", err)
	}
	span.End()
	temperature := models.NewTemperature(city, tempC)

	return &temperature, nil
}

func validateZipcode(zipcode string) error {
	if len(zipcode) != 8 {
		return ErrInvalidZipcode
	}

	pattern := `^\d{8}$`
	match, _ := regexp.MatchString(pattern, zipcode)
	if !match {
		return ErrInvalidZipcode
	}

	return nil
}
