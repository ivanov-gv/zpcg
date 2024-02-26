# Montenegro Railways Timetable Bot

A Telegram bot that simplifies finding train schedules in Montenegro.

**[@Monterails_bot](https://t.me/Monterails_bot)**

# Why

Original timetable - https://zpcg.me/search

The site does not have a mobile version, pretty hard to navigate and can't show intersection routes.
This bot is an attempt to solve those issues.

## Requirements

1. Ensure it's straightforward and user-friendly to the greatest extent possible
   1. ✅ Provides a menu with all the supported commands
   2. ✅ Has a start message with laconic and complete description of how to get the timetable
   3. ~ Returns the requested timetable even if the user made a typo
   4. ~ In case of an error returns message with a detailed description of how to resolve the issue
   5. ✅ Has enough logging to check if the users reach their goals
   6. ❌ Recognizes latin and cyrillic alphabet
   7. ❌ Provides button interface for the most used stations
   8. ❌ Provides additional information:
      1. ❌ Price and discounts
      2. ❌ Marks trains as local/fast
      3. ❌ Stations location, ticket office availability
      4. ❌ Official website links
      5. ❌ Last update time
2. Has a full, updated timetable
   1. ✅ Parses original website automatically
   2. ✅ Knows every single station
   3. ❌ Updates timetable automatically
   4. ❌ Updates timetable only if needed (checks if there are any changes)
3. ✅ Cost-effective as much as possible
   1. ✅ Runs in the cloud
   2. ✅ Uses tg webhooks
   3. ✅ Scales up/down automatically
   4. ✅ Runs on 1~2 thread cpus with minimum RAM available on the cloud

## Known issues

# Solution details

## Interface

### /start

### Timetable request

#### No intersections

#### With an intersection

### Error message

## Backend

### Path finding algorithm

#### Assumptions

#### Future-proof

#### Steps

### Parser

#### Parsing algo

#### Timetable storage

### Telegram bot

#### Start-up and scalability

#### Approximate match for the station names

#### Path finder

#### Render

### DevOps

# Ways to improve

## Interface

### Keyboard

### Inline buttons

### Links

### Cyrillic alphabet

### Telegram WebApp

## Additional info

## Other platforms