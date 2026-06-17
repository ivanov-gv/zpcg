# Montenegro Railways Timetable Bot

Stateless Telegram bot for Montenegro Railways timetable lookup. Timetable compiled into the binary at build time — no
DB, no cache, no runtime deps. Domain-specific O(n) path-finding on a tree-topology railway network. ~€0.03/month on
Cloud Run.

Deployed at **[@Monterails_bot](https://t.me/Monterails_bot)** since January 2024, actively maintained.
Serving ~150 unique users per month and ~500 requests per month. Breaking changes are not expected - the scope is
intentionally fixed.

## Documentation

Start with:

- [Overview](docs/1-the-problem-overview.md) - for the reasoning behind the project and solution architecture.
- [System Design](docs/2-system-design.md) - for the detailed design of the bot's architecture and its components.

Or any of the following:

| Doc                                                           | Covers                                                                                            |
|---------------------------------------------------------------|---------------------------------------------------------------------------------------------------|
| [1. Overview](docs/1-the-problem-overview.md)                 | The problem, the solution, traffic estimation, scope, and constraints                             |
| [2. System Design](docs/2-system-design.md)                   | Functional and non-functional requirements, high-level architecture, capacity estimates           |
| [3. Implementation Details](docs/3-implementation-details.md) | Bot interface, the path-finding algorithm and its correctness assumptions, backend design         |
| [4. Development Guidelines](docs/4-development-guidelines.md) | Coding standards and conventions, shared across human and AI contributors                         |
| [5. Operational Profile](docs/5-operational-profile.md)       | CI/CD, environments, release and rollback, observability, failure recovery, real production costs |
| [6. Extension Roadmap](docs/6-execution-roadmap.md)           | Scoped future extensions and the trade-offs of each                                               |
