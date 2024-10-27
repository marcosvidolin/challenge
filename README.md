# Challenge

## Run

1. To start all the infra and the api run:

```bash
make up
```

2. To start the producer run

```
make produce
```

## References

for more details see the documentation:

[API Documentation](api/README.md)

[Producer Documentation](producer/README.md)

## Time Frame

| Date    | Hours | Project                  |
|---------|-------|--------------------------|
| 23/Oct  | 3h    | producer                 |
| 24/Oct  | 2h    | producer                 |
| 25/Oct  | 3h    | infra / docker-compose   |
| 26/Oct  | 8h    | api                      |
| 27/Oct  | 6h    | api                      |

Total: 22h

## Test Difficult

Challenging

Its not very complex to develop, but there are some rules to check and a lot of work to do within a
short time frame. The hardest part was managing time. I had to discard some things early, like unit
tests and proper documentation. Later I noticed certain details in the CSV items that would take
significant work to handle.
The main issues I found were related to IDs (see the session "Known Issues" for more details).
Overall, it was a good challenge.

## Known Issues

- "Foreign key constraint": This occurs when attempting to insert a user record that references
  a foreign key which does not exist yet. To address this I could catch this error and re-queue the
  user to be processed again. If the error persists after retrying, the user record should be sent
  to a Dead Letter Queue (DLQ) for further investigation or manual processing.

- "ID reference": Im referencing an ID from the CSV file.

- "Upsert Item by Item": I can perform the upsert with many users at time, this will be good
  for not open many connections to the database
