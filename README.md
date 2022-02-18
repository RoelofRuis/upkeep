# Upkeep

A console application for time tracking.

Install locally with `make install` or build an executable with `make build-prod`.

## Usage

Upkeep tracks time in blocks on a per-day basis. Each block belongs to a category. You can view totals per day or category and create exports of time ranges.

Run `upkeep help` to view the available commands.

Most commands accept a date (range) parameter given as `d:<date>`.
Here `<date>` can be one of the following:
- Any date in the form of `YYYY-MM-DD` to select that specific date.
- `d` or `day` to select a day.
- `w` or `week` to select a week, monday to friday.
- `wf` or `weekfull` to select a week, monday to sunday.
- `m` or `month` to select a month.

The letter/word inputs can be prefixed with (negative) numbers to shift their time range. For instance to select two weeks ago, use `d:-2w`.

#### Data storage

Configuration and past blocks are stored as JSON files in the `.upkeep` folder in the user home directory.

## Development

Build development executable with `make build` or a debug executable with `make build-dbg`. 