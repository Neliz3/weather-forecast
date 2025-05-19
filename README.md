# Weather Forecast Service

Weather RESTful API application that allows users to subscribe to weather updates for their city via email. It uses [weatherapi.com](https://www.weatherapi.com/) API for weather data and <b>Mailchimp</b> for email subscription management.

## Tools
* Gin
* PostgreSQL
* golang-migrate
* Docker
* testing (httptest)
* Render (deploy)
* Weather API 
* SendGrid (email)
* golangci-lint (linter)
* JWT

## TDD Philosophy
The application follows the Test-Driven Development (TDD) approach to ensure robust and reliable code.

## API Endpoints
The API provides the following endpoints as defined in [`swagger.yml`](swagger.yml):

### Weather
- **GET /weather**: Fetch the current weather for a specified city.

### Subscription
- **POST /subscribe**: Subscribe an email to receive weather updates for a specific city.
- **GET /confirm/{token}**: Confirm an email subscription using a token.
- **GET /unsubscribe/{token}**: Unsubscribe from weather updates using a token.

## Getting Started
1. Clone the repository.
2. Set up the [`.env`](.env) file with your API keys and secrets (refer to [`example.env`](example.env) for structure).
3. Build and start the services with Docker Compose
```
docker compose up --build
```
4. To run migrations, use
```
docker compose up -d db
docker compose run --rm migrate
docker compose build app
```

## Future Updates
* Endpoint's testing and mock servers
* Add achive emails logic (don't delete permanently)
* Decode key query parameter in /weather endpoint with JWT token
