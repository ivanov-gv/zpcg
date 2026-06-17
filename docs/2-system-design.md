# System Design

See [1-overview.md](1-overview.md) for the reasons behind this design.

## Functional requirements

1. It's a Telegram bot
2. Provides a menu with all available commands
3. Simple commands for providing general information:
    1. `/start` - a start message with a brief and complete description of how to get the timetable
    2. `/help` - frequently asked questions
    3. `/map` - a railway map
    4. `/about` - a brief description of the bot and the link to the source code
4. A user sends a message with a timetable request: "Podgorica, Bar"
    1. The bot must recognise both Latin and Cyrillic alphabet
    2. If the input is valid, i.e. contains exactly two stations (departure, arrival) - the bot responds with an
       up-to-date timetable
    3. If the input contains some minor typos, the bot corrects them and responds with an up-to-date timetable
    4. If the requested station is not found or does not exist - the bot responds with an error message
    5. If the input is invalid - the bot responds with an error message
5. The response with a timetable
    1. Contains the departure and arrival times for all trains on the route, sorted by the departure time
    2. Shows the last update time
    3. Has a link to the timetable on the official website

## Non-functional requirements

### Availability

High availability, low rate of failure. Even if the zpcg.me website is down, the bot must still work.

### Scalability

As stated in the [1-overview.md](1-overview.md) file, we expect low traffic, no more than 100 requests per day.
The bot must be able to scale down to zero and scale up automatically on demand.

### Low latency

Response time should be around 300ms on average with a maximum of 1 second.

### Durability

We do not expect the bot to store any user data. All observability metrics should be stored for no less than 1 month.

### Read to write ratio

The timetable changes only 3 times a year, while requests are serviceable every day.

## High-level design

### API

Only one API endpoint for receiving requests from Telegram webhooks.

```http request
GET /
X-Telegram-Bot-Api-Secret-Token: # (optional) a predefined secret token. See http://core.telegram.org/bots/api
content-type: application/json

{
    "update_id": integer,
    "message": jsonObject,          # if a user sends a message. see https://core.telegram.org/bots/api#message
    "callback_query": jsonObject    # if a user presses a callback button. see https://core.telegram.org/bots/api#callbackquery
}

response:
200 OK
500 Internal Server Error
```

### Data modeling

No persistent data storage, fully stateless architecture. All data needed for timetable is parsed from the official site
and compiled into the binary.

### Estimates




