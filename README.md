# ZJUT-class-timetable-ics
A conversion script for zjuter from zf system class timetable to iCal.

Call WeJH's Api (login, get the class table)
And Transfer the raw json data to icalendar.

# Features

- [x] Generate `.ics` file.
- [ ] Convert Exam Events
- [ ]

# How to use

1. clone the repo
2. `cp config.example.yml config.yml` and edit it.
3. `go run .`

and one file named `class.ics` will be created.

# Others

PRs are welcome.
