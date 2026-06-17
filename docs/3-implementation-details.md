# Implementation

## Interface

### Commands (functional requirements 2. and 3.)

As stated in the [requirements (2. and 3.)](2-system-design.md#functional-requirements) section, the bot must have a
number of slash
commands.

#### /start

`/start` is the first command received when a user activates the bot. It returns quick,
laconic instructions on how to use the bot.

Returns the message in different languages depending on the user's Telegram settings.
If the user's language is not supported, it uses English for communication.

#### /help

Returns FAQ

#### /map

Prints a link to the map with all stations: https://goo.gl/maps/jYRwCAC14mmoe4L3A

#### /about

Gives links to this project and the licence description

### Timetable request (functional requirements 4.)

Any message without a leading `/` is parsed as a timetable request.

Let's introduce two definitions:

`{delimiter symbol}` - one of !@#$%^&*()_+-=[]{}|;':",./<>? chars

`{departure/arrival station name}` - any number of other chars

The request is expected to be in the following formats:

`{departure station name} {delimiter symbol} {arrival station name}`

or

`{departure station name} {delimiter symbol} {any number of any chars} {delimiter symbol} {arrival station name}`

The format is built this way to be typos-proof as required by
[the system design doc (4.3)](2-system-design.md#functional-requirements). If a user types an input message:

- with several delimiters in a row - the bot responds without an error
- with a typo in one of two station names - the bot is able to find the closest match more precisely, since stations
  names are easily dividable by the delimiter symbol

Messages in the expected format receive a timetable. For any other format, there will be an error.

The default error is defined in `/internal/service/render/{language_code}/*.go` files.

If the input message is in the right format but contains a station name from the blacklist defined in
`/internal/service/blacklist` - it is responded to with a custom warning usually telling the user
that the station does not exist as required
by [the system design doc (4.4)](2-system-design.md#functional-requirements).

### Timetable response (functional requirements 5.)

According to the [system design doc (5.)](2-system-design.md#functional-requirements):

> 5 - The response with a timetable
>   1. Contains the departure and arrival times for all trains on the route, sorted by the departure time
>   2. Shows the last update time
>   3. Has a link to the timetable on the official website

#### No intersections

Response with a route without an intersection is in the following format:

```
        Podgorica > Danilovgrad
[7100](...) 08:00 > 08:29 
[7102](...) 12:50 > 13:19 
[7104](...) 15:35 > 16:04 
[7106](...) 18:30 > 18:59 
[7108](...) 21:45 > 22:14

{a link to the route on zpcg.me}

{update button with the last update time} {reverse route button}
```

1. Header contains two station names (departure > arrival) written as in the official timetable
2. Headers' delimiter character '>' is aligned to match the timetable character if possible
3. The following rows contain a train number with a link to the official timetable, departure time from the
   departure station, and departure time from the arrival station
4. Rows are sorted in ascending order by the departure station time
5. The entire message uses a monospace font to make it possible to match the indent of the header and
   the timetable rows

#### With an intersection

A response for a route with an intersection follows this format:

```
Virpazar > Podgorica > Nikšić
[6100](...) 05:43 > 06:19 
[6150](...) 06:55 > 07:26 
[7100](...)         08:00 > 09:03 
[6152](...) 09:51 > 10:22 
[6154](...) 12:08 > 12:39 
[7102](...)         12:50 > 13:53 
[6156](...) 14:41 > 15:12 
[7104](...)         15:35 > 16:38 
[6102](...) 16:03 > 17:06 
[6158](...) 17:35 > 18:06 
[7106](...)         18:30 > 19:33 
[6104](...) 18:25 > 19:15 
[7108](...)         21:45 > 22:48 
[6160](...) 21:29 > 22:00

{a link to the route on zpcg.me}

{update button with the last update time} {reverse route button}
```

The format is the same as for the route without intersections but:

1. There are three columns for the times: departure time from the departure station,
   departure time from the intersection station, and the departure time from the arrival station
2. Rows are sorted by the time of the intersection station

### Error message

Default:

```
Try again - two stations, separated by a comma. Just like that:

Podgorica, Niksic
```

Or a custom one for blacklisted stations:

```
There is no train station in Budva
```

## Backend

Language choice is Go (Golang), since it is a

- compiled language
- With a strong community and a large ecosystem of libraries and tools. This includes the '%#v' verb in
  the https://pkg.go.dev/fmt package, perfectly suitable for code generation. Timetable compile-time injection is
  made possible with this option - it's a requirement from [the system design doc](2-system-design.md#data-modeling)
- With a strong focus on performance - response time [constraint](2-system-design.md#latency) of ~300 ms is achievable
  with this choice

and the last, but not least - it's the most familiar language to the dev of the project.

[Non-functional requirements](2-system-design.md#non-functional-requirements) explicitly state that the backend must not
have any dependencies on external services and all necessary data must be pre-parsed and compiled into the binary.
That's why we need Exporter - a tool for fetching the data from zpcg.me API and exporting it as a Go code file ready
to be compiled.

### Exporter

Exporter: `cmd/exporter`

The Exporter prepares the necessary data structures for finding routes. It parses
the [official site](https://zpcg.me/search)
and saves these structures as a Go code file in `gen/timetable/timetable.gen.go`. The entire timetable is then compiled
directly into the executable file, making it accessible without any additional steps.

### Path-finding algorithm

The easiest way to solve the problem is to use the Dijkstra's algorithm with O(n^2) complexity.
However, a number of assumptions that are true for the Montenegro railway system allow us to build a more efficient
algorithm **with linear complexity**.

It allows us to meet the [latency](2-system-design.md#latency) requirements and
reduce [costs](2-system-design.md#costs).

#### Assumptions

![img.png](docs/resources/img.png)

The railway system of Montenegro consists of:

- Main station Podgorica
- Podgorica - Bar branch
- Podgorica - Niksic branch
- Podgorica - Belgrade branch (through Bijelo Polje)
- Podgorica - Tuzi (and further to Durrës/Tirane, Albania. This branch is fully abandoned and is not in use)

So the assumptions are:

1. Each train passes every station in its route only once.
2. For any two stations there is a straight route or a route with an intersection in Podgorica station.
   This means the Podgorica station might be used as the only intersection station. Also, there are no routes
   with two or more intersections.
3. Routes without intersections are always faster (or just preferred) than those with one or more.
4. Routes with an intersection in Podgorica are always preferred to the routes with any other.
5. The Montenegro railway system is an acyclic undirected graph. This means that for every train station,
   trains can travel in any direction, and it is impossible to circle from a station back to itself
   without repeating stations.

This means that it is completely enough to return only one of two types of routes for every possible input:

1. Straight route
2. Route with exactly one interchange in Podgorica

#### Future-proof

Conditions that might disrupt the assumptions include:

- If the timetable includes a train not passing through the Podgorica station
- If the railway system has a route with any other preferred interchange station other than Podgorica
- If a cycle occurs in the graph

Possible improvements to this railway system might involve:

Extending the Podgorica - Belgrade branch to the Subotica station and beyond to Hungary, in light of Subotica
being mentioned in the [official timetable](https://www.zcg-prevoz.me/search) and the Subotica - Szeged rail line is now
[reopened](https://www.railwaypro.com/wp/subotica-szeged-rail-passenger-services-resumed/). It may add new stations
to the branch but is not going to create a cycle. Some kind of fork might be added which means there might be another
preferred intersection point for trains going from Hungary to Serbia, but it is obviously outside the context of
the Montenegro timetable bot.

The same applies to the possible renovation of the Podgorica - Tuzi - Durrës/Tirane branch and an extension of
the Podgorica - Niksic line to the Sarajevo.

It means we can easily rely on the assumptions listed above to optimize the path-finding algorithm.

#### Steps

1. **Preparation: needed data structures**
    1. **StationIdToTrainIdSetMap**: map: station -> set of train ids

       This field is a map where each key is a StationId and each value is a set of TrainId.
       This map allows the PathFinder to know which trains depart from any given station. This information is essential
       for
       identifying possible routes when calculating the paths between two stations.

       Train ids are unique and for every station exists a set of trains that pass the station. Consequently, it is
       possible to fill this structure correctly.

    2. **TrainIdToStationsMap**: map: train id -> (set of stations -> {station name, arrival time, departure time})

       This field is also a map. Each key is a TrainId and each value is a StationIdToStationMap
       that contains every station that the train stops at. This map is important to derive the sequence of stations
       that
       each train traverses along its route.

       According to the assumptions, every train passes each station on its route only once. That makes it possible to
       build a set of route stations with details about the arrival/departure time.

    3. **TransferStation**: station id for the interchange station

       This field represents a predefined station which serves as the only transfer point in case there are no direct
       paths available between the two requested stations.

2. **Check for Direct Path**

   Find a direct path/route from station A (aStation) to station B (bStation).
   "Direct" here means there are trains that go from station A to station B without the need for a transfer.

   It is done by first checking which trains serve both station A and station B using the
   predefined StationIdToTrainIdSetMap:
   ```go
       trainIdSetA = p.stationIdToTrainIdSetMap[aStation]
       trainIdSetB = p.stationIdToTrainIdSetMap[bStation]
       // get the intersection of maps of the trains 
       possibleRoutes := utils.Intersection(trainIdSetA, trainIdSetB)
   ```
   If the set intersection is not empty - the direct paths are found.

   The final result is defined by validating if the train/s found has a schedule such that it/they depart/s
   from station A and arrive at station B consecutively. This information might be found using prepared
   TrainIdToStationsMap. If the condition is met, then this direct path is valid and is added to the returning paths.

   We assume that the direct paths are the preferred ones, so there is no need to look for other paths.

   The problem is solved.

3. **Check for Indirect Path**

   If no direct path is found, paths with a transfer need to be found. This is the step where the predefined
   Transfer Station is being used.

   In this step the algorithm essentially performs a two-part journey - first, it finds the paths from station A
   to the Transfer Station, then it finds the paths from the Transfer Station to station B the same way as it was done
   in the previous step. These possible routes are merged in a specific manner, in case of overlapping
   schedules,
   the route which arrives at the transfer station before the departure to station B is preferred.

   We assume that the path through the Transfer station always exists and is the most suitable.

4. **Returning Paths**

   Finally, return the paths found and a bool flag indicating whether the
   found routes are direct or not.

#### Concurrency-safety

Every step of the algorithm:

- Uses read operations from maps and slices
- Does not modify any of the shared structures, no write operations

Which means that the algorithm is safe to run concurrently without any additional synchronisation.

#### Complexity

In the first step, the algorithm iterates over each train in a set that contains the intersection of trains
at station A and station B. The intersection operation would have a time complexity of O(n), where n is the number of
elements in the larger train set. Then, for every train, it checks the conditions and possibly appends it to the paths.
Adding all of this up, the time complexity of this operation can be considered to be O(n).

In the second one, the algorithm repeats step one twice, so that would be 2*O(n) ~ O(n). Then it iterates through the
lists
of direct paths from station A to the transfer station and from the transfer station to station B, merging these paths
as it goes. This merger can, in the worst case, have a time complexity of O(n + m), where n and m are the lengths of the
two lists. Adding these together, the overall time complexity could be considered to be O(n + m) where n is the number
of elements in the larger set and m is the total number of paths.

The space complexity of the algorithm comes primarily from the storage of the 'paths'. In the worst case,
it could potentially store each direct path from the source station to the destination, as well as each path
via the transfer station. This gives a worst-case space complexity of O(n + m), where n and m are the numbers of
direct and transfer paths, respectively.

Please note: In these analyses, m and n are not necessarily numbers of stations or trains but rather the number
of paths and matches the algorithm needs to keep track of. Thus, the actual time or space this algorithm takes might be
different based on the specific graph structure of the train network.

Also, the memory requirement of the StationIdToTrainIdSetMap and TrainIdToStationsMap structures scale linearly with
the overall increase in the number of stations and trains. The exact growth can be viewed as O(n) and O(m) respectively,
where 'n' is the number of train services and 'm' is the number of stations.

Therefore, the algorithm has **a linear time and space complexity.**

### Telegram bot serverless app

Initializer: `cmd/tg-init`
Server: `cmd/tg-server`

Cloud Run was selected because its request-based billing model matches a 100%-read, low-traffic workload exactly: zero
idle cost, and the free tier covers the estimated load entirely (see [Estimates](2-system-design.md#estimates)).
Familiarity with the platform is also a key factor in the choice - GCloud and Cloud Run are both familiar to the
dev of the project.

The bot waits for HTTP requests from Telegram, handles, and sends messages to users. It runs on the Google Cloud Run
platform, invokes on every request (i.e. is not running all time long) which makes it cost-effective.

The life cycle is simple:

1. A user request is received. If there are no running instances of the bot - a new one is started. It takes less than a
   second to wake up and respond correctly to a live probe.
2. Request is handled, a route is generated, and a reply message is sent. It takes no more than 200ms.
3. If there are no more requests for ~15 minutes - Cloud Run automatically scales it all down

Overall request handling time is 200ms if a request is processed in a warm state. If a cold startup is required, then
a 1-second delay is added, resulting in a 1200ms response time in the worst-case. These numbers fit perfectly into
the [latency](2-system-design.md#latency) requirements.

The request handling time (~200ms) is the only billable time in this lifecycle. The overall cost of the bot right now
is about 0.03 euros per month with 100-150 users and ~500 user requests per 30 days. The required 1€ per month
[budget](2-system-design.md#costs) is more than enough for the project.

The bot is almost free to maintain and as cost-effective as possible.
